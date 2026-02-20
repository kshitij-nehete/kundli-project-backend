package domain

import "time"

type User struct {
	ID           string    `bson:"_id,omitempty"`
	Name         string    `bson:"name"`
	Email        string    `bson:"email"`
	PasswordHash string    `bson:"password_hash"`
	IsPremium    bool      `bson:"is_premium"`
	ReportCount  int       `bson:"report_count"`
	CreatedAt    time.Time `bson:"created_at"`
}
