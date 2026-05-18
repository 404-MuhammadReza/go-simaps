package repository

import (
	"time"
	"context"

	"go-simaps/internal/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type LogRepository interface {
	Create(ctx context.Context, log model.UsageLog) error
	GetSummary(ctx context.Context, today, week, month, sixMonths, year time.Time) ([]model.Summary, error)
	GetDetails(ctx context.Context, feature string, start, end time.Time) ([]model.UsageLog, error)
}

type logRepository struct { collection *mongo.Collection }

func NewLogRepository(db *mongo.Database) LogRepository {
	collectionName := "usage_logs"

	collections, err := db.ListCollectionNames(context.Background(), bson.M{"name": collectionName})
	if err == nil && len(collections) == 0 {
		tsOptions := options.CreateCollection().
			SetTimeSeriesOptions(options.TimeSeries().
				SetTimeField("timestamp").
				SetMetaField("metadata").
				SetGranularity("minutes"))

		_ = db.CreateCollection(context.Background(), collectionName, tsOptions)
	}

	return &logRepository{ collection: db.Collection(collectionName) }
}

func (r *logRepository) Create(ctx context.Context, log model.UsageLog) error {
	_, err := r.collection.InsertOne(ctx, log)
	if err != nil { return err }
	return nil
}

func (r *logRepository) GetSummary(ctx context.Context, today, week, month, sixMonths, year time.Time) ([]model.Summary, error) {
	matchStage := bson.D{{Key: "$match", Value: bson.M{"timestamp": bson.M{"$gte": year}}}}
	groupStage := bson.D{{
		Key: "$group",
		Value: bson.M{
			"_id":        "$feature",
			"today":      bson.M{"$sum": bson.M{"$cond": bson.A{bson.M{"$gte": bson.A{"$timestamp", today}}, 1, 0}}},
			"week":       bson.M{"$sum": bson.M{"$cond": bson.A{bson.M{"$gte": bson.A{"$timestamp", week}}, 1, 0}}},
			"month":      bson.M{"$sum": bson.M{"$cond": bson.A{bson.M{"$gte": bson.A{"$timestamp", month}}, 1, 0}}},
			"six_months": bson.M{"$sum": bson.M{"$cond": bson.A{bson.M{"$gte": bson.A{"$timestamp", sixMonths}}, 1, 0}}},
			"year":       bson.M{"$sum": 1},
		},
	}}

	sortStage := bson.D{{Key: "$sort", Value: bson.D{{Key: "_id", Value: 1}}}}
	cursor, err := r.collection.Aggregate(ctx, mongo.Pipeline{matchStage, groupStage, sortStage})
	if err != nil { return nil, err }
	defer cursor.Close(ctx)

	recaps := make([]model.Summary, 0)
	err = cursor.All(ctx, &recaps); 
	if err != nil { return nil, err }

	return recaps, nil
}

func (r *logRepository) GetDetails(ctx context.Context, feature string, start, end time.Time) ([]model.UsageLog, error) {
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{
			"feature":   feature,
			"timestamp": bson.M{"$gte": start, "$lte": end},
		}}},
		{{Key: "$lookup", Value: bson.M{
			"from":         "users",
			"localField":   "user_id",
			"foreignField": "_id",
			"as":           "user",
		}}},
		{{Key: "$addFields", Value: bson.M{
			"user_name": bson.M{"$arrayElemAt": bson.A{"$user.name", 0}},
		}}},
		{{Key: "$sort", Value: bson.D{{Key: "timestamp", Value: -1}}}},
	}

	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil { return nil, err }
	defer cursor.Close(ctx)

	logs := make([]model.UsageLog, 0)
	err = cursor.All(ctx, &logs)
	if err != nil { return nil, err }

	return logs, nil
}