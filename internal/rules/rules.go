package rules

import (
	"strings"
	"time"

	"crowdreview/internal/models"
)

// RuleResult expresses the effect of a single fraud rule.
type RuleResult struct {
	Name     string
	Passed   bool
	Score    float64
	Severity string
	Details  map[string]interface{}
}

// RunAll evaluates all rules and returns their results.
func RunAll(review models.Review) []RuleResult {
	return []RuleResult{
		textLengthRule(review),
		extremeRatingRule(review),
		suspiciousLanguageRule(review),
		geoRule(review),
		freshAccountRule(review),
		ipFrequencyRule(review),
	}
}

func textLengthRule(review models.Review) RuleResult {
	words := len(strings.Fields(review.Content))
	score := 10.0
	passed := words >= 20
	if !passed {
		score = -15.0
	}
	return RuleResult{
		Name:     "text_length",
		Passed:   passed,
		Score:    score,
		Severity: ternary(passed, "low", "medium"),
		Details: map[string]interface{}{
			"word_count": words,
		},
	}
}

func extremeRatingRule(review models.Review) RuleResult {
	passed := review.Rating != 1 && review.Rating != 5
	score := 5.0
	if !passed {
		score = -10.0
	}
	return RuleResult{
		Name:     "rating_discrepancy",
		Passed:   passed,
		Score:    score,
		Severity: ternary(passed, "low", "medium"),
		Details: map[string]interface{}{
			"rating": review.Rating,
		},
	}
}

func suspiciousLanguageRule(review models.Review) RuleResult {
	lower := strings.ToLower(review.Content)
	badWords := []string{"free money", "click here", "guaranteed", "fake", "scam"}
	passed := true
	for _, b := range badWords {
		if strings.Contains(lower, b) {
			passed = false
			break
		}
	}
	score := 8.0
	if !passed {
		score = -25.0
	}
	return RuleResult{
		Name:     "language_filter",
		Passed:   passed,
		Score:    score,
		Severity: ternary(passed, "low", "high"),
		Details: map[string]interface{}{
			"matched": !passed,
		},
	}
}

func geoRule(review models.Review) RuleResult {
	passed := review.GeoLocation != "" && review.GeoLocation != "unknown"
	score := 4.0
	if !passed {
		score = -5.0
	}
	return RuleResult{
		Name:     "geolocation",
		Passed:   passed,
		Score:    score,
		Severity: ternary(passed, "low", "low"),
		Details: map[string]interface{}{
			"geo": review.GeoLocation,
		},
	}
}

func freshAccountRule(review models.Review) RuleResult {
	// We lack the user creation timestamp here; rely on CreatedAt meta if present.
	created := review.CreatedAt
	fresh := time.Since(created) < 24*time.Hour
	score := 6.0
	passed := !fresh
	if fresh {
		score = -12.0
	}
	return RuleResult{
		Name:     "fresh_account",
		Passed:   passed,
		Score:    score,
		Severity: ternary(passed, "low", "medium"),
		Details: map[string]interface{}{
			"fresh": fresh,
		},
	}
}

func ipFrequencyRule(review models.Review) RuleResult {
	// Placeholder: in a real system we'd check Redis counts; here we flag empty IP.
	passed := review.IPAddress != ""
	score := 5.0
	if !passed {
		score = -10.0
	}
	return RuleResult{
		Name:     "ip_presence",
		Passed:   passed,
		Score:    score,
		Severity: ternary(passed, "low", "medium"),
		Details: map[string]interface{}{
			"ip": review.IPAddress,
		},
	}
}

func ternary[T any](cond bool, a, b T) T {
	if cond {
		return a
	}
	return b
}
