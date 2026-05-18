package service

import (
	"context"

	"go-simaps/internal/model"
	"go-simaps/internal/apperror"
	"go-simaps/internal/repository"
)

type UserService interface {
	GetAll(ctx context.Context) ([]model.UserResponse, error)
	GetByID(ctx context.Context, userID string) (*model.UserResponse, error)
	UpdateRole(ctx context.Context, userID string, request model.RoleRequest) (*model.UserResponse, error)
	Delete(ctx context.Context, userID string) error
}

type userService struct {
	repository repository.UserRepository
}

func NewUserService(repository repository.UserRepository) UserService {
	return &userService{repository}
}

func (s *userService) GetAll(ctx context.Context) ([]model.UserResponse, error) {
	users, err := s.repository.FindAll(ctx)
	if err != nil { return nil, apperror.ErrInternal("Failed to retrieve users") }

	response := make([]model.UserResponse, 0, len(users))
	for i := range users {
		response = append(response, users[i].ToResponse())
	}

	return response, nil
}

func (s *userService) GetByID(ctx context.Context, userID string) (*model.UserResponse, error) {
	user, err := s.repository.FindByID(ctx, userID)
	if err != nil { return nil, apperror.ErrInternal("Failed to retrieve user") }
	if user == nil { return nil, apperror.ErrNotFound("User not found") }

	response := user.ToResponse()
	return &response, nil
}

func (s *userService) UpdateRole(ctx context.Context, userID string, request model.RoleRequest) (*model.UserResponse, error) {
	user, err := s.repository.FindByID(ctx, userID)
	if err != nil { return nil, apperror.ErrInternal("Failed to find user") }
	if user == nil { return nil, apperror.ErrNotFound("User not found") }

	user.PatchFrom(&request)
	err = s.repository.Update(ctx, *user)
	if err != nil { return nil, apperror.ErrInternal("Failed to update user") }

	response := user.ToResponse()
	return &response, nil
}

func (s *userService) Delete(ctx context.Context, userID string) error {
	exist, err := s.repository.ExistsByID(ctx, userID)
	if err != nil { return apperror.ErrInternal("Failed to check user existence") }
	if !exist { return apperror.ErrNotFound("User not found") }

	err = s.repository.Delete(ctx, userID)
	if err != nil { return apperror.ErrInternal("Failed to delete user") }
	return nil
}
