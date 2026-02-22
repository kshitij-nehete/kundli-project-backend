package repository

import (
	"context"
	"time"

	"github.com/kshitij-nehete/astro-report/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoReportRepository struct {
	collection *mongo.Collection
}

func NewMongoReportRepository(db *mongo.Database) *MongoReportRepository {
	return &MongoReportRepository{
		collection: db.Collection("reports"),
	}
}

func (r *MongoReportRepository) Create(ctx context.Context, report *domain.Report) error {

	report.CreatedAt = time.Now()
	report.ExpiresAt = report.CreatedAt.Add(90 * 24 * time.Hour)

	_, err := r.collection.InsertOne(ctx, report)
	return err
}

func (r *MongoReportRepository) CountByUser(ctx context.Context, userID string) (int64, error) {

	return r.collection.CountDocuments(ctx, bson.M{
		"user_id": userID,
	})
}
