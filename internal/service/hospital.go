package service

import (
	"io"
	"time"
	"bytes"
	"context"
	"strings"

	"go-simaps/internal/utils"
	"go-simaps/internal/model"
	"go-simaps/internal/constant"
	"go-simaps/internal/apperror"
	"go-simaps/internal/repository"
)

type HospitalService interface {
	GetAll(ctx context.Context) ([]model.HospitalResponse, error)
	Create(ctx context.Context, userID string, request model.HospitalRequest) (*model.HospitalResponse, error)
	Update(ctx context.Context, userID string, hospitalID string, request model.HospitalRequest) (*model.HospitalResponse, error)
	Delete(ctx context.Context, hospitalID string) error

	CreateMany(ctx context.Context, userID string, file io.Reader) (*model.HospitalBulkResponse, error)
	GetReport(payload model.HospitalBulkResponse) (*bytes.Buffer, error)
	GetTemplate() (*bytes.Buffer, error)
	Export(ctx context.Context) (*bytes.Buffer, error)
}

type hospitalService struct {
	repository repository.HospitalRepository
	geocoding  GeocodingService
	logging    LoggingService
}

func NewHospitalService(repository repository.HospitalRepository, geocoding GeocodingService, logging LoggingService) HospitalService {
	return &hospitalService{repository, geocoding, logging}
}

func (s *hospitalService) GetAll(ctx context.Context) ([]model.HospitalResponse, error) {
	hospitals, err := s.repository.FindAll(ctx)
	if err != nil { return nil, apperror.ErrInternal("Failed to retrieve hospitals") }

	response := make([]model.HospitalResponse, 0, len(hospitals))
	for i := range hospitals {
		response = append(response, hospitals[i].ToResponse())
	}

	return response, nil
}

func (s *hospitalService) Create(ctx context.Context, userID string, request model.HospitalRequest) (*model.HospitalResponse, error) {
	hospitalID := utils.GenerateID(constant.HospitalPrefix)
	coordinate := request.Coordinate

	if coordinate == nil {
		response, err := s.geocoding.ToCoordinate(ctx, request.Address)
		if err != nil { return nil, err }
		coordinate = response

		s.logging.LogUsage(ctx, userID, model.LogRequest{
			Feature: "Geocoding API",
			Action:  "Create Hospital: " + request.Name,
		})
	}

	hospital := request.ToModel(hospitalID, *coordinate)
	err := s.repository.Create(ctx, hospital)
	if err != nil { return nil, apperror.ErrInternal("Failed to create hospital") }

	response := hospital.ToResponse()
	return &response, nil
}

func (s *hospitalService) Update(ctx context.Context, userID string, hospitalID string, request model.HospitalRequest) (*model.HospitalResponse, error) {
	hospital, err := s.repository.FindByID(ctx, hospitalID)
	if err != nil { return nil, apperror.ErrInternal("Failed to find hospital") }
	if hospital == nil { return nil, apperror.ErrNotFound("Hospital not found") }

	coordinate  := hospital.Coordinate
	if request.Coordinate != nil && *request.Coordinate != hospital.Coordinate {
		coordinate = *request.Coordinate
	} else if request.Address != hospital.Address {
		response, err := s.geocoding.ToCoordinate(ctx, request.Address)
		if err != nil { return nil, err }
		coordinate = *response

		s.logging.LogUsage(ctx, userID, model.LogRequest{
			Feature: "Geocoding API",
			Action:  "Update Hospital: " + hospital.Name,
		})
	}

	hospital.PutFrom(request, coordinate)
	err = s.repository.Update(ctx, *hospital)
	if err != nil { return nil, apperror.ErrInternal("Failed to update hospital") }

	response := hospital.ToResponse()
	return &response, nil
}

func (s *hospitalService) Delete(ctx context.Context, hospitalID string) error {
	exist, err := s.repository.ExistsByID(ctx, hospitalID)
	if err != nil { return apperror.ErrInternal("Failed to check hospital existence") }
	if !exist { return apperror.ErrNotFound("Hospital not found") }

	err = s.repository.Delete(ctx, hospitalID)
	if err != nil { return apperror.ErrInternal("Failed to delete hospital") }
	return nil
}

func (s *hospitalService) CreateMany(ctx context.Context, userID string, file io.Reader) (*model.HospitalBulkResponse, error) {
	requests, err := utils.ParseHospitalExcel(file)
	if err != nil { return nil, apperror.ErrBadRequest(err.Error()) }

	var success []model.HospitalResponse
	var failed []model.HospitalFailedDetail
	for i, request := range requests {
		errMsgs := utils.ValidateHospitalExcel(request)
		if len(errMsgs) > 0 {
			detail := model.CreateHospitalFailedDetail(i+1, errMsgs, &request)
			failed = append(failed, detail)
			continue
		}

		response, err := s.Create(ctx, userID, request)
		if err != nil {
			detail := model.CreateHospitalFailedDetail(i+1, []string{err.Error()}, &request)
			failed = append(failed, detail)
			continue
		}

		success = append(success, *response)
		time.Sleep(100 * time.Millisecond)
	}

	response := model.CreateHospitalBulkResponse(success, failed)
	return &response, nil
}

func (s *hospitalService) GetReport(payload model.HospitalBulkResponse) (*bytes.Buffer, error) {
	var success [][]interface{}
	for _, data := range payload.SuccessData {
		success = append(success, []interface{}{
			data.Name, data.Phone, data.Address, data.Coordinate.Lat,
			data.Coordinate.Lng, "SUCCESS", "",
		})
	}

	var failed [][]interface{}
	for _, data := range payload.FailedData {
		var lat, lng float64
		if data.Value.Coordinate != nil {
			lat = data.Value.Coordinate.Lat
			lng = data.Value.Coordinate.Lng
		}

		errMsg := strings.Join(data.Error, ", ")
		failed = append(failed, []interface{}{
			data.Value.Name, data.Value.Phone, data.Value.Address,
			lat, lng, "FAILED", errMsg,
		})
	}

	columns := utils.HospitalReportColumns
	excel, err := utils.GenerateReport(columns, success, failed)
	if err != nil { return nil, apperror.ErrInternal("Failed to generate Excel") }

	return excel, nil
}

func (s *hospitalService) GetTemplate() (*bytes.Buffer, error) {
	columns := utils.HospitalColumns
	excel, err := utils.GenerateTemplate(columns)
	if err != nil { return nil, apperror.ErrInternal("Failed to generate Excel") }
	return excel, nil
}

func (s *hospitalService) Export(ctx context.Context) (*bytes.Buffer, error) {
	hospitals, err := s.repository.FindAll(ctx)
	if err != nil { return nil, apperror.ErrInternal("Failed to retrieve hospitals") }

	var data [][]interface{}
	for _, hospital := range hospitals {
		data = append(data, []interface{}{
			hospital.Name, hospital.Phone, hospital.Address,
			hospital.Coordinate.Lat, hospital.Coordinate.Lng,
		})
	}

	columns := utils.HospitalColumns
	excel, err := utils.GenerateExport(columns, data)
	if err != nil { return nil, apperror.ErrInternal("Failed to generate Excel") }

	return excel, nil
}
