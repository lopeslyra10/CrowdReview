package validation

import (
	"context"
	"log"
	"time"

	"crowdreview/internal/models"
	"crowdreview/internal/repository"
)

// FraudWorker consumes the fraud-validation-queue asynchronously.
type FraudWorker struct {
	Queue      chan models.Review
	Engine     *FraudEngine
	Validation repository.ValidationRepository
}

// FraudQueueName provides a friendly identifier for observability/logs.
const FraudQueueName = "fraud-validation-queue"

func NewFraudWorker(engine *FraudEngine, validationRepo repository.ValidationRepository) *FraudWorker {
	return &FraudWorker{
		Queue:      make(chan models.Review, 100),
		Engine:     engine,
		Validation: validationRepo,
	}
}

// Enqueue pushes a review for validation without blocking the request lifecycle.
func (w *FraudWorker) Enqueue(review models.Review) {
	select {
	case w.Queue <- review:
	default:
		log.Println("fraud-validation-queue is full, dropping review", review.ID)
	}
}

// Start begins processing the queue. Should run in a goroutine.
func (w *FraudWorker) Start() {
	go func() {
		for review := range w.Queue {
			w.process(review)
		}
	}()
}

func (w *FraudWorker) process(review models.Review) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, suspicious := w.Engine.Evaluate(review)
	if err := w.Validation.SaveResult(ctx, &result); err != nil {
		log.Printf("failed to save validation result: %v", err)
		return
	}

	status := "approved"
	if suspicious {
		status = "flagged"
	}
	if err := w.Validation.MarkReview(ctx, review.ID, result.ID, status, suspicious); err != nil {
		log.Printf("failed to mark review: %v", err)
	}
}
