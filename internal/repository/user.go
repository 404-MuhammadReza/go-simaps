package repository

import (
	"context"

	"go-simaps/internal/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository interface {
	FindAll(ctx context.Context) ([]model.User, error)
	FindByID(ctx context.Context, userID string) (*model.User, error)
	FindByMicrosoftID(ctx context.Context, microsoftId string) (*model.User, error)
	ExistsByID(ctx context.Context, userID string) (bool, error)

	Create(ctx context.Context, user model.User) error
	Update(ctx context.Context, user model.User) error
	Delete(ctx context.Context, userID string) error
}

type userRepository struct { collection *mongo.Collection }

func NewUserRepository(db *mongo.Database) UserRepository {
	return &userRepository{ collection: db.Collection("users") }
}

func (r *userRepository) FindAll(ctx context.Context) ([]model.User, error) {
	var users []model.User
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil { return nil, err }
	defer cursor.Close(ctx)

	err = cursor.All(ctx, &users)
	if err != nil { return nil, err }
	return users, nil
}

func (r *userRepository) FindByID(ctx context.Context, userID string) (*model.User, error) {
	var user model.User

	filter := bson.M{"_id": userID}
	err := r.collection.FindOne(ctx, filter).Decode(&user)
	if err == mongo.ErrNoDocuments { return nil, nil }
	if err != nil { return nil, err }

	return &user, nil
}

func (r *userRepository) FindByMicrosoftID(ctx context.Context, microsoftId string) (*model.User, error) {
	var user model.User

	filter := bson.M{"microsoft_id": microsoftId}
	err := r.collection.FindOne(ctx, filter).Decode(&user)
	if err == mongo.ErrNoDocuments { return nil, nil }
	if err != nil { return nil, err }

	return &user, nil
}

func (r *userRepository) ExistsByID(ctx context.Context, userID string) (bool, error) {
	filter := bson.M{"_id": userID}
	count, err := r.collection.CountDocuments(ctx, filter)
	if err != nil { return false, err }

	return count > 0, nil
}

func (r *userRepository) Create(ctx context.Context, user model.User) error {
	_, err := r.collection.InsertOne(ctx, user)
	if err != nil { return err }

	return nil
}

func (r *userRepository) Update(ctx context.Context, user model.User) error {
	filter := bson.M{"_id": user.ID}
	update := bson.M{"$set": user}
	_, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil { return err }

	return nil
}

func (r *userRepository) Delete(ctx context.Context, userID string) error {
	filter := bson.M{"_id": userID}
	_, err := r.collection.DeleteOne(ctx, filter)
	if err != nil { return err }

	return nil
}