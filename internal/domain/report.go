package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ReportStatus string

const (
	StatusProcessing ReportStatus = "processing"
	StatusActive     ReportStatus = "active"
	StatusExpired    ReportStatus = "expired"
)

type Report struct {
	ID              primitive.ObjectID `bson:"_id,omitempty"`
	UserID          string             `bson:"user_id"`
	Name            string             `bson:"name"`
	BirthDate       string             `bson:"birth_date"`
	BirthTime       string             `bson:"birth_time"`
	PlaceOfBirth    string             `bson:"place_of_birth,omitempty"`
	PlanetaryData   interface{}        `bson:"planetary_data,omitempty"`
	AIReport        interface{}        `bson:"ai_report,omitempty"`
	ConfidenceScore int                `bson:"confidence_score,omitempty"`
	Status          ReportStatus       `bson:"status"`
	CreatedAt       time.Time          `bson:"created_at"`
	ExpiresAt       time.Time          `bson:"expires_at"`
	NumerologyData  *NumerologyReport  `bson:"numerology_data,omitempty"`
}

func (r *Report) IsExpired() bool {
	return time.Now().After(r.ExpiresAt)
}
