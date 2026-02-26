package ai

import (
	"errors"

	"github.com/kshitij-nehete/astro-report/internal/domain"
)

func ValidateNumerologyReport(r domain.NumerologyReport) error {

	if len(r.PersonalityOutlook) < 5 || len(r.PersonalityOutlook) > 10 {
		return errors.New("invalid personality_outlook count")
	}

	if len(r.CareerPrediction) != 5 {
		return errors.New("career_prediction must have 5 items")
	}

	if len(r.WealthPrediction) != 10 {
		return errors.New("wealth_prediction must have 10 items")
	}

	if len(r.MarriageRelationship) != 10 {
		return errors.New("marriage_relationship must have 10 items")
	}

	if len(r.HealthPrediction) != 5 {
		return errors.New("health_prediction must have 5 items")
	}

	if len(r.Remedies) != 5 {
		return errors.New("remedies must have 5 items")
	}

	return nil
}
