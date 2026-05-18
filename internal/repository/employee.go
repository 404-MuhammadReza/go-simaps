package repository

import (
	"context"
	
	"go-simaps/internal/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type EmployeeRepository interface {
	FindAll(ctx context.Context) ([]model.Employee, error)
	FindByID(ctx context.Context, employeeID string) (*model.Employee, error)
	FindByEmail(ctx context.Context, email string) (*model.Employee, error)

	ExistsByID(ctx context.Context, employeeID string) (bool, error)
	ExistsByEmployeeID(ctx context.Context, employeeID string) (bool, error)
	ExistsByEmail(ctx context.Context, email string) (bool, error)

	Create(ctx context.Context, employee model.Employee) error
	Update(ctx context.Context, employee model.Employee) error
	Delete(ctx context.Context, employeeID string) error
	DeletePicture(ctx context.Context, employeeID string) error
}

type employeeRepository struct { collection *mongo.Collection }
func NewEmployeeRepository(db *mongo.Database) EmployeeRepository {
	return &employeeRepository{ collection: db.Collection("employees") }
}

func (r *employeeRepository) FindAll(ctx context.Context) ([]model.Employee, error) {
	var employees []model.Employee

	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil { return nil, err }
	defer cursor.Close(ctx)

	err = cursor.All(ctx, &employees)
	if err != nil { return nil, err }
	return employees, nil
}

func (r *employeeRepository) FindByID(ctx context.Context, employeeID string) (*model.Employee, error) {
	var employee model.Employee

	filter := bson.M{"_id": employeeID}
	err := r.collection.FindOne(ctx, filter).Decode(&employee)
	if err == mongo.ErrNoDocuments { return nil, nil }
	if err != nil { return nil, err }

	return &employee, nil
}

func (r *employeeRepository) FindByEmail(ctx context.Context, email string) (*model.Employee, error) {
	var employee model.Employee

	filter := bson.M{"email": email}
	err := r.collection.FindOne(ctx, filter).Decode(&employee)
	if err == mongo.ErrNoDocuments { return nil, nil }
	if err != nil { return nil, err }

	return &employee, nil
}

func (r *employeeRepository) ExistsByID(ctx context.Context, employeeID string) (bool, error) {
	filter := bson.M{"_id": employeeID}
	count, err := r.collection.CountDocuments(ctx, filter)
	if err != nil { return false, err }

	return count > 0, nil
}

func (r *employeeRepository) ExistsByEmployeeID(ctx context.Context, employeeID string) (bool, error) {
	filter := bson.M{"employee_id": employeeID}
	count, err := r.collection.CountDocuments(ctx, filter)
	if err != nil { return false, err }

	return count > 0, nil
}

func (r *employeeRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	filter := bson.M{"email": email}
	count, err := r.collection.CountDocuments(ctx, filter)
	if err != nil { return false, err }

	return count > 0, nil
}

func (r *employeeRepository) Create(ctx context.Context, employee model.Employee) error {
	_, err := r.collection.InsertOne(ctx, employee)
	if err != nil { return err }

	return nil
}

func (r *employeeRepository) Update(ctx context.Context, employee model.Employee) error {
	filter := bson.M{"_id": employee.ID}
	_, err := r.collection.ReplaceOne(ctx, filter, employee)
	if err != nil { return err }

	return nil
}

func (r *employeeRepository) Delete(ctx context.Context, employeeID string) error {
	filter := bson.M{"_id": employeeID}
	_, err := r.collection.DeleteOne(ctx, filter)
	if err != nil { return err }

	return nil
}

func (r *employeeRepository) DeletePicture(ctx context.Context, employeeID string) error {
	filter := bson.M{"_id": employeeID}
	update := bson.M{"$unset": bson.M{"picture_object": ""}}

	_, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil { return err }

	return nil
}