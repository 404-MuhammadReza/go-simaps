package repository

import (
	"context"
	
	"go-simaps/internal/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type HospitalRepository interface {
	FindAll(ctx context.Context) ([]model.Hospital, error)
	FindByID(ctx context.Context, hospitalID string) (*model.Hospital, error)
	ExistsByID(ctx context.Context, hospitalID string) (bool, error)
	Create(ctx context.Context, hospital model.Hospital) error
	Update(ctx context.Context, hospital model.Hospital) error
	Delete(ctx context.Context, hospitalID string) error
}

type hospitalRepository struct { collection *mongo.Collection }
func NewHospitalRepository(db *mongo.Database) HospitalRepository {
	return &hospitalRepository{ collection: db.Collection("hospitals") }
}

func (r *hospitalRepository) FindAll(ctx context.Context) ([]model.Hospital, error) {
	var hospitals []model.Hospital

	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil { return nil, err }
	defer cursor.Close(ctx)

	err = cursor.All(ctx, &hospitals)
	if err != nil { return nil, err }
	return hospitals, nil
}

func (r *hospitalRepository) FindByID(ctx context.Context, hospitalID string) (*model.Hospital, error) {
	var hospital model.Hospital

	filter := bson.M{"_id": hospitalID}
	err := r.collection.FindOne(ctx, filter).Decode(&hospital)
	if err == mongo.ErrNoDocuments { return nil, nil }
	if err != nil { return nil, err }

	return &hospital, nil
}

func (r *hospitalRepository) ExistsByID(ctx context.Context, hospitalID string) (bool, error) {
	filter := bson.M{"_id": hospitalID}
	count, err := r.collection.CountDocuments(ctx, filter)
	if err != nil { return false, err }

	return count > 0, nil
}

func (r *hospitalRepository) Create(ctx context.Context, hospital model.Hospital) error {
	_, err := r.collection.InsertOne(ctx, hospital)
	if err != nil { return err }

	return nil
}

func (r *hospitalRepository) Update(ctx context.Context, hospital model.Hospital) error {
	filter := bson.M{"_id": hospital.ID}
	_, err := r.collection.ReplaceOne(ctx, filter, hospital)
	if err != nil { return err }

	return nil
}

func (r *hospitalRepository) Delete(ctx context.Context, hospitalID string) error {
	filter := bson.M{"_id": hospitalID}
	_, err := r.collection.DeleteOne(ctx, filter)
	if err != nil { return err }

	return nil
}