package model

import (
	"time"

	"go-simaps/internal/constant"
)

// User Model
type User struct {
	ID 			string	  `bson:"_id"`
	Name 		string	  `bson:"name"`
	Email 		string	  `bson:"email"`
	MicrosoftID string	  `bson:"microsoft_id"`
	Role 		string	  `bson:"role"`

	CreatedAt	time.Time `bson:"created_at"`
	UpdatedAt	time.Time `bson:"updated_at"`
}

// DTOs
// Role Request DTO
type RoleRequest struct {
	Role string `json:"role" binding:"required,oneof=admin super_admin user"`
}

// User Response DTO
type UserResponse struct {
	ID    		 string    `json:"id"`
	Name  		 string    `json:"name"`
	Email 		 string    `json:"email"`
	MicrosoftID  string    `json:"microsoft_id"`
	Role 		 string    `json:"role"`

	CreatedAt	 time.Time `json:"created_at"`
	UpdatedAt	 time.Time `json:"updated_at"`
}

// Azure AD Response DTO
type AzureResponse struct {
	ID				  string `json:"id"`
	DisplayName		  string `json:"displayName"`
	Mail			  string `json:"mail"`
	UserPrincipalName string `json:"userPrincipalName"`
}

// Mapping Functions
// Azure AD Response to Model
func (azure *AzureResponse) ToModel(userID string) User {
	return User{
		ID: userID,
		Name: azure.DisplayName,
		Email: azure.Mail,
		MicrosoftID: azure.ID,
		Role: constant.User,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// Model User to Response
func (user *User) ToResponse() UserResponse {
	return UserResponse{
		ID: user.ID,
		Name: user.Name,
		Email: user.Email,
		MicrosoftID: user.MicrosoftID,
		Role: user.Role,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

// Sync User Model with Azure AD Response
func (user *User) SyncWithAzure(azure *AzureResponse) {
	user.Name = azure.DisplayName
	user.Email = azure.Mail
}

// Patch update for User model
func (user *User) PatchFrom(request *RoleRequest) {
	user.Role = request.Role
	user.UpdatedAt = time.Now()
}