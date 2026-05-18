package service

import (
	"context"
	"encoding/json"

	"go-simaps/internal/utils"
	"go-simaps/internal/model"
	"go-simaps/internal/config"
	"go-simaps/internal/constant"
	"go-simaps/internal/apperror"
	"go-simaps/internal/repository"

	"golang.org/x/oauth2"
)

type AuthService interface {
	GenerateLoginURL() (string, string, error)
	ExchangeToken(ctx context.Context, code string) (string, string, error)
	RefreshToken(ctx context.Context, refreshToken string) (string, error)
	DevSession(ctx context.Context, userID string) (string, error)
}

type authService struct {
	repository 		 repository.UserRepository
	oauth      		 *oauth2.Config
	JWTSecret  		 []byte
	JWTIssuer  		 string
	JWTAccessExpiry  int
	JWTRefreshExpiry int
}

func NewAuthService(repository repository.UserRepository, config *config.Config, oauth *oauth2.Config) AuthService {
	return &authService{
		repository: repository,
		oauth: oauth,
		JWTSecret: config.JWTSecret,
		JWTIssuer: config.JWTIssuer,
		JWTAccessExpiry: config.JWTAccessExpiry,
		JWTRefreshExpiry: config.JWTRefreshExpiry,
	}
}

func (s *authService) GenerateLoginURL() (string, string, error) {
	state, err := utils.GenerateState()
	if err != nil { return "", "", apperror.ErrInternal("Failed to generate state") }

	authOptions := []oauth2.AuthCodeOption{
		oauth2.SetAuthURLParam("prompt", "select_account"),
		oauth2.SetAuthURLParam("max_age", "0"),
	}

	url := s.oauth.AuthCodeURL(state, authOptions...)
	return url, state, nil
}

func (s *authService) ExchangeToken(ctx context.Context, code string) (string, string, error) {
	response, err := s.fetchAzure(ctx, code)
	if err != nil { return "", "", err }

	user, err := s.upsertUser(ctx, response)
	if err != nil { return "", "", err }

	accessToken, err := utils.GenerateAccessToken(user.ID, user.Email, user.Role, s.JWTIssuer, s.JWTSecret, s.JWTAccessExpiry)
	if err != nil { return "", "", apperror.ErrInternal("Failed to generate access tokens") }

	refreshToken, err := utils.GenerateRefreshToken(user.ID, s.JWTIssuer, s.JWTSecret, s.JWTRefreshExpiry)
	if err != nil { return "", "", apperror.ErrInternal("Failed to generate refresh tokens") }

	return accessToken, refreshToken, nil
}

func (s *authService) RefreshToken(ctx context.Context, refreshToken string) (string, error) {
	claims, err := utils.ValidateRefreshToken(refreshToken, s.JWTSecret)
	if err != nil { return "", apperror.ErrUnauthorized("Invalid refresh token") }

	user, err := s.repository.FindByID(ctx, claims.UserID)
	if err != nil { return "", apperror.ErrInternal("Failed to find user") }
	if user == nil { return "", apperror.ErrNotFound("User not found") }

	accessToken, err := utils.GenerateAccessToken(user.ID, user.Email, user.Role, s.JWTIssuer, s.JWTSecret, s.JWTAccessExpiry)
	if err != nil { return "", apperror.ErrInternal("Failed to generate access token") }

	return accessToken, nil
}

func (s *authService) DevSession(ctx context.Context, userID string) (string, error) {
	user, err := s.repository.FindByID(ctx, userID)
	if err != nil { return "", apperror.ErrInternal("Failed to find user") }
	if user == nil { return "", apperror.ErrNotFound("User not found") }

	accessToken, err := utils.GenerateAccessToken(user.ID, user.Email, user.Role, s.JWTIssuer, s.JWTSecret, s.JWTAccessExpiry)
	if err != nil { return "", apperror.ErrInternal("Failed to generate access token") }

	return accessToken, nil
}

// Helper functions for token exchange flow
// Fetch user info from Microsoft Graph API using the authorization code
func (s *authService) fetchAzure(ctx context.Context, code string) (*model.AzureResponse, error) {
	microsoftToken, err := s.oauth.Exchange(ctx, code)
	if err != nil { return nil, apperror.ErrExternal("Failed to exchange token") }

	client := s.oauth.Client(ctx, microsoftToken)
	response, err := client.Get("https://graph.microsoft.com/v1.0/me")
	if err != nil { return nil, apperror.ErrExternal("Failed to get user info") }
	defer response.Body.Close()

	var azureResponse model.AzureResponse
	err = json.NewDecoder(response.Body).Decode(&azureResponse)
	if err != nil { return nil, apperror.ErrExternal("Failed to decode user info") }

	return &azureResponse, nil
}

// Upsert user data based on Azure response (create new or update existing)
func (s *authService) upsertUser(ctx context.Context, azure *model.AzureResponse) (*model.User, error) {
	user, err := s.repository.FindByMicrosoftID(ctx, azure.ID)
	if err != nil { return nil, apperror.ErrInternal("Failed to check user existence") }

	if user != nil {
		user.SyncWithAzure(azure)
		err = s.repository.Update(ctx, *user)
		if err != nil { return nil, apperror.ErrInternal("Failed to update user") }

		return user, nil
	}

	userID := utils.GenerateID(constant.UserPrefix)
	newUser := azure.ToModel(userID)
	err = s.repository.Create(ctx, newUser)
	if err != nil { return nil, apperror.ErrInternal("Failed to create user") }

	return &newUser, nil
}