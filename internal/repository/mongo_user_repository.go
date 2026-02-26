package repository

import (
	"context"
	"errors"
	"time"

	"github.com/kshitij-nehete/astro-report/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoUserRepository struct {
	collection *mongo.Collection
}

func NewMongoUserRepository(db *mongo.Database) *MongoUserRepository {
	return &MongoUserRepository{
		collection: db.Collection("users"),
	}
}

func (r *MongoUserRepository) Create(ctx context.Context, user *domain.User) error {

	user.CreatedAt = time.Now()

	_, err := r.collection.InsertOne(ctx, user)
	return err
}

func (r *MongoUserRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {

	var user domain.User

	err := r.collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &user, nil
}

func (r *MongoReportRepository) UpdateStatus(
	ctx context.Context,
	reportID string,
	status domain.ReportStatus,
) error {

	objectID, err := primitive.ObjectIDFromHex(reportID)
	if err != nil {
		return err
	}

	_, err = r.collection.UpdateByID(
		ctx,
		objectID,
		bson.M{"$set": bson.M{"status": status}},
	)
	return err
}

func (r *MongoReportRepository) FindByID(
	ctx context.Context,
	reportID string,
) (*domain.Report, error) {

	objectID, err := primitive.ObjectIDFromHex(reportID)
	if err != nil {
		return nil, err
	}

	var report domain.Report

	err = r.collection.FindOne(
		ctx,
		bson.M{"_id": objectID},
	).Decode(&report)

	if err != nil {
		return nil, err
	}

	return &report, nil
}

func (r *MongoReportRepository) FindByUser(
	ctx context.Context,
	userID string,
	limit int64,
) ([]*domain.Report, error) {

	opts := options.Find().
		SetSort(bson.M{"created_at": -1}).
		SetLimit(limit)

	cursor, err := r.collection.Find(
		ctx,
		bson.M{"user_id": userID},
		opts,
	)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var reports []*domain.Report
	if err := cursor.All(ctx, &reports); err != nil {
		return nil, err
	}

	return reports, nil
}
