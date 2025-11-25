package validation

import (
	"crowdreview/internal/models"
	"crowdreview/internal/rules"
)

// FraudEngine aggregates rule scores into a final confidence metric.
type FraudEngine struct{}

func NewFraudEngine() *FraudEngine {
	return &FraudEngine{}
}

// Evaluate runs all rules and returns a validation result populated with signals.
func (f *FraudEngine) Evaluate(review models.Review) (models.ReviewValidationResult, bool) {
	results := rules.RunAll(review)

	score := 50.0
	signals := make([]models.FraudSignal, 0, len(results))
	checks := make(map[string]interface{})

	for _, res := range results {
		score += res.Score
		checks[res.Name] = res.Details
		if !res.Passed {
			signals = append(signals, models.FraudSignal{
				Type:     res.Name,
				Severity: res.Severity,
				Details:  res.Details,
			})
		}
	}

	if score < 0 {
		score = 0
	}
	if score > 100 {
		score = 100
	}

	outcome := "approved"
	suspicious := false
	if score < 55 {
		outcome = "flagged"
		suspicious = true
	} else if score < 40 {
		outcome = "rejected"
		suspicious = true
	}

	return models.ReviewValidationResult{
		ReviewID: review.ID,
		Score:    score,
		Outcome:  outcome,
		Checks:   checks,
		Signals:  signals,
	}, suspicious
}
