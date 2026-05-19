package service

import (
	"time"
	"bytes"
	"context"

	"go-simaps/internal/utils"
	"go-simaps/internal/model"
	"go-simaps/internal/apperror"
	"go-simaps/internal/repository"
)

type LoggingService interface {
	LogUsage(ctx context.Context, userID string, request model.LogRequest) error
	GetSummary(ctx context.Context) ([]model.Summary, error)
	GetDetails(ctx context.Context, query model.DetailRequest) ([]model.UsageLog, error)
	ExportDetails(ctx context.Context, query model.DetailRequest) (*bytes.Buffer, error)
}

type loggingService struct {
	repository repository.LogRepository
}

func NewLoggingService(repository repository.LogRepository) LoggingService {
	return &loggingService{repository}
}

func (s *loggingService) LogUsage(ctx context.Context, userID string, request model.LogRequest) error {
	model := request.ToModel(userID)
	err := s.repository.Create(ctx, model)
	if err != nil { return apperror.ErrInternal("Failed to log usage") }

	return nil
}

func (s *loggingService) GetSummary(ctx context.Context) ([]model.Summary, error) {
	now := time.Now()

	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	week := today.AddDate(0, 0, -7)
	month := today.AddDate(0, -1, 0)
	sixMonths := today.AddDate(0, -6, 0)
	year := today.AddDate(-1, 0, 0)

	response, err := s.repository.GetSummary(ctx, today, week, month, sixMonths, year)
	if err != nil { return nil, apperror.ErrInternal("Failed to retrieve usage recap") }

	return response, nil
}

func (s *loggingService) GetDetails(ctx context.Context, query model.DetailRequest) ([]model.UsageLog, error) {
	start, end, err := parseDetailRange(query)
	if err != nil { return nil, err }

	return s.repository.GetDetails(ctx, query.Feature, start, end)
}

func (s *loggingService) ExportDetails(ctx context.Context, query model.DetailRequest) (*bytes.Buffer, error) {
	start, end, err := parseDetailRange(query)
	if err != nil { return nil, err }

	logs, err := s.repository.GetDetails(ctx, query.Feature, start, end)
	if err != nil { return nil, apperror.ErrInternal("Failed to retrieve usage details") }

	data := make([][]interface{}, 0, len(logs))
	for _, log := range logs {
		data = append(data, []interface{}{
			log.Timestamp.Format("2006-01-02 15:04:05"),
			*log.UserName, log.Feature, log.Action,
		})
	}

	columns := utils.LogDetailColumns
	excel, err := utils.GenerateExport(columns, data)
	if err != nil { return nil, apperror.ErrInternal("Failed to generate Excel") }

	return excel, nil
}

func parseDetailRange(query model.DetailRequest) (time.Time, time.Time, error) {
	layout := "2006-01-02"
	start, err := time.Parse(layout, query.StartDate)
	if err != nil { return time.Time{}, time.Time{}, apperror.ErrBadRequest("invalid start_date format, use YYYY-MM-DD") }

	end, err := time.Parse(layout, query.EndDate)
	if err != nil { return time.Time{}, time.Time{}, apperror.ErrBadRequest("invalid end_date format, use YYYY-MM-DD") }

	end = time.Date(end.Year(), end.Month(), end.Day(), 23, 59, 59, 999999999, end.Location())
	if start.After(end) { return time.Time{}, time.Time{}, apperror.ErrBadRequest("start_date cannot be after end_date") }

	return start, end, nil
}
