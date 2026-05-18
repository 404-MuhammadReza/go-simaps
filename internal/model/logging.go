package model

import "time"

// UsageLog Model
type UsageLog struct {
    Timestamp time.Time	`bson:"timestamp" json:"timestamp"`
	UserID    string	`bson:"user_id" json:"user_id"`
    Feature   string	`bson:"feature" json:"feature"`
    Action    string	`bson:"action" json:"action"`
	
	UserName  *string	`bson:"user_name,omitempty" json:"user_name,omitempty"`
}

// Summary Model
type Summary struct {
	Feature   string `bson:"_id" json:"feature"`
	Today     int64  `bson:"today" json:"today"`
	Week      int64  `bson:"week" json:"week"`
	Month     int64  `bson:"month" json:"month"`
	SixMonths int64  `bson:"six_months" json:"six_months"`
	Year      int64  `bson:"year" json:"year"`
}

// DTOs
// DetailUsage Request DTO
type DetailRequest struct {
	Feature   string `form:"feature" binding:"required"`
	StartDate string `form:"start_date" binding:"required"`
	EndDate   string `form:"end_date" binding:"required"`
}

// LogRequest DTOs
type LogRequest struct {
	Feature string `json:"feature" binding:"required"`
	Action  string `json:"action" binding:"required"`
}

// Mapping function for request to model
func (request *LogRequest) ToModel(userID string) UsageLog {
	return UsageLog{
		Timestamp: time.Now(),
		Feature: request.Feature,
		UserID: userID,
		Action: request.Action,
	}
}