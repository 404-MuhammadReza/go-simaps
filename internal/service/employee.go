package service

import (
	"io"
	"fmt"
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

type EmployeeService interface {
	GetAll(ctx context.Context) ([]model.EmployeeResponse, error)
	GetByEmail(ctx context.Context, userEmail string) (*model.EmployeeResponse, error)

	Create(ctx context.Context, userID string, request model.EmployeeRequest) (*model.EmployeeResponse, error)
	Update(ctx context.Context, userID string, employeeID string, request model.EmployeeRequest) (*model.EmployeeResponse, error)
	Delete(ctx context.Context, employeeID string) error

	UpdatePicture(ctx context.Context, employeeID string, picture io.Reader, size int64, contentType string) (*model.EmployeeResponse, error)
	DeletePicture(ctx context.Context, employeeID string) (*model.EmployeeResponse, error)

	CreateMany(ctx context.Context, userID string, file io.Reader) (*model.EmployeeBulkResponse, error)
	GetReport(payload model.EmployeeBulkResponse) (*bytes.Buffer, error)
	GetTemplate() (*bytes.Buffer, error)
	Export(ctx context.Context) (*bytes.Buffer, error)
}

type employeeService struct {
	repository repository.EmployeeRepository
	storage	   repository.StorageRepository
	geocoding  GeocodingService
	logging    LoggingService
}

func NewEmployeeService(repository repository.EmployeeRepository, storage repository.StorageRepository, geocoding GeocodingService, logging LoggingService) EmployeeService {
	return &employeeService{repository, storage, geocoding, logging}
}

func (s *employeeService) GetAll(ctx context.Context) ([]model.EmployeeResponse, error) {
	employees, err := s.repository.FindAll(ctx)
	if err != nil { return nil, apperror.ErrInternal("Failed to retrieve employees") }

	response := make([]model.EmployeeResponse, 0, len(employees))
	for i := range employees {
		var url string
		if employees[i].Picture != "" {
			url, _ = s.storage.GetPresignedURL(ctx, employees[i].Picture, time.Hour)
		}
		response = append(response, employees[i].ToResponse(url))
	}

	return response, nil
}

func (s *employeeService) GetByEmail (ctx context.Context, userEmail string) (*model.EmployeeResponse, error) {
	employee, err := s.repository.FindByEmail(ctx, userEmail)
	if err != nil { return nil, apperror.ErrInternal("Failed to retrieve employee") }
	if employee == nil { return nil, apperror.ErrNotFound("Employee not found") }

	var url string
	if employee.Picture != "" {
		url, _ = s.storage.GetPresignedURL(ctx, employee.Picture, time.Hour)
	}

	response := employee.ToResponse(url)
	return &response, nil
}

func (s *employeeService) Create(ctx context.Context, userID string, request model.EmployeeRequest) (*model.EmployeeResponse, error) {
	existByID, err := s.repository.ExistsByEmployeeID(ctx, request.EmployeeID)
	if err != nil { return nil, apperror.ErrInternal("Failed to check employee existence") }
	if existByID { return nil, apperror.ErrConflict("Employee with the same ID already exists") }

	existByEmail, err := s.repository.ExistsByEmail(ctx, request.Email)
	if err != nil { return nil, apperror.ErrInternal("Failed to check employee existence") }
	if existByEmail { return nil, apperror.ErrConflict("Employee with the same email already exists") }

	employeeID := utils.GenerateID(constant.EmployeePrefix)
	coordPrimary := request.CoordPrimary
	if coordPrimary == nil {
		response, err := s.geocoding.ToCoordinate(ctx, request.AddressPrimary)
		if err != nil { return nil, err }
		coordPrimary = response

		s.logging.LogUsage(ctx, userID, model.LogRequest{
			Feature: "Geocoding API",
			Action:  "Create Employee (Primary): " + request.Name,
		})
	}

	var coordSecondary *model.Coordinate
	if request.AddressSecondary != "" {
		if request.CoordSecondary != nil {
			coordSecondary = request.CoordSecondary
		} else {
			response, err := s.geocoding.ToCoordinate(ctx, request.AddressSecondary)
			if err != nil { return nil, err }
			coordSecondary = response

			s.logging.LogUsage(ctx, userID, model.LogRequest{
				Feature: "Geocoding API",
				Action:  "Create Employee (Secondary): " + request.Name,
			})
		}
	}

	employee := request.ToModel(employeeID, *coordPrimary, coordSecondary)
	err = s.repository.Create(ctx, employee)
	if err != nil { return nil, apperror.ErrInternal("Failed to create employee") }

	response := employee.ToResponse("")
	return &response, nil
}

func (s *employeeService) Update(ctx context.Context, userID string, employeeID string, request model.EmployeeRequest) (*model.EmployeeResponse, error) {
	employee, err := s.repository.FindByID(ctx, employeeID)
	if err != nil { return nil, apperror.ErrInternal("Failed to find employee") }
	if employee == nil { return nil, apperror.ErrNotFound("Employee not found") }

	if request.EmployeeID != "" && request.EmployeeID != employee.EmployeeID {
		existByID, err := s.repository.ExistsByEmployeeID(ctx, request.EmployeeID)
		if err != nil { return nil, apperror.ErrInternal("Failed to check employee existence") }
		if existByID { return nil, apperror.ErrConflict("Employee with the same ID already exists") }
	}

	if request.Email != "" && request.Email != employee.Email {
		existByEmail, err := s.repository.ExistsByEmail(ctx, request.Email)
		if err != nil { return nil, apperror.ErrInternal("Failed to check employee existence") }
		if existByEmail { return nil, apperror.ErrConflict("Employee with the same email already exists") }
	}

	coordPrimary := employee.CoordPrimary
	if request.CoordPrimary != nil && *request.CoordPrimary != employee.CoordPrimary {
		coordPrimary = *request.CoordPrimary
	} else if request.AddressPrimary != employee.AddressPrimary {
		response, err := s.geocoding.ToCoordinate(ctx, request.AddressPrimary)
		if err != nil { return nil, err }
		coordPrimary = *response

		s.logging.LogUsage(ctx, userID, model.LogRequest{
			Feature: "Geocoding API",
			Action:  "Update Employee (Primary): " + employee.Name,
		})
	}

	coordSecondary := employee.CoordSecondary
	if request.CoordSecondary != nil && (employee.CoordSecondary == nil || *request.CoordSecondary != *employee.CoordSecondary) {
		coordSecondary = request.CoordSecondary
	} else if request.AddressSecondary != "" && request.AddressSecondary != employee.AddressSecondary {
		response, err := s.geocoding.ToCoordinate(ctx, request.AddressSecondary)
		if err != nil { return nil, err }
		coordSecondary = response

		s.logging.LogUsage(ctx, userID, model.LogRequest{
			Feature: "Geocoding API",
			Action:  "Update Employee (Secondary): " + employee.Name,
		})
	} else if request.AddressSecondary == "" { coordSecondary = nil }

	employee.PutFrom(request, coordPrimary, coordSecondary)
	err = s.repository.Update(ctx, *employee)
	if err != nil { return nil, apperror.ErrInternal("Failed to update employee") }

	var url string
	if employee.Picture != "" { url, _ = s.storage.GetPresignedURL(ctx, employee.Picture, time.Hour) }

	response := employee.ToResponse(url)
	return &response, nil
}

func (s *employeeService) Delete(ctx context.Context, employeeID string) error {
	employee, err := s.repository.FindByID(ctx, employeeID)
	if err != nil { return apperror.ErrInternal("Failed to find employee") }
	if employee == nil { return apperror.ErrNotFound("Employee not found") }

	err = s.repository.Delete(ctx, employeeID)
	if err != nil { return apperror.ErrInternal("Failed to delete employee") }
	if employee.Picture != "" { _ = s.storage.Delete(ctx, employee.Picture) }

	return nil
}

func (s *employeeService) UpdatePicture(ctx context.Context, employeeID string, picture io.Reader, size int64, contentType string) (*model.EmployeeResponse, error) {
	employee, err := s.repository.FindByID(ctx, employeeID)
	if err != nil { return nil, apperror.ErrInternal("Failed to find employee") }
	if employee == nil { return nil, apperror.ErrNotFound("Employee not found") }
	if employee.Picture != "" { _ = s.storage.Delete(ctx, employee.Picture) }

	extension := contentType[strings.LastIndex(contentType, "/")+1:]
	objectName := fmt.Sprintf("pict_%s.%s", employee.EmployeeID, extension)
	_, err = s.storage.Upload(ctx, objectName, picture, size, contentType)
	if err != nil { return nil, apperror.ErrInternal("Failed to upload employee picture") }

	employee.Picture = objectName
	err = s.repository.Update(ctx, *employee)
	if err != nil { return nil, apperror.ErrInternal("Failed to update employee") }

	url, _ := s.storage.GetPresignedURL(ctx, objectName, time.Hour)
	response := employee.ToResponse(url)
	return &response, nil
}

func (s *employeeService) DeletePicture(ctx context.Context, employeeID string) (*model.EmployeeResponse, error) {
	employee, err := s.repository.FindByID(ctx, employeeID)
	if err != nil { return nil, apperror.ErrInternal("Failed to find employee") }
	if employee == nil { return nil, apperror.ErrNotFound("Employee not found") }

	if employee.Picture != "" {
		err = s.storage.Delete(ctx, employee.Picture)
		if err != nil { return nil, apperror.ErrInternal("Failed to delete employee picture") }
	}

	err = s.repository.DeletePicture(ctx, employeeID)
	if err != nil { return nil, apperror.ErrInternal("Failed to update employee") }

	response := employee.ToResponse("")
	return &response, nil
}

func (s *employeeService) CreateMany(ctx context.Context, userID string, file io.Reader) (*model.EmployeeBulkResponse, error) {
	requests, err := utils.ParseEmployeeExcel(file)
	if err != nil { return nil, apperror.ErrBadRequest(err.Error()) }

	var success []model.EmployeeResponse
	var failed []model.EmployeeFailedDetail
	for i, request := range requests {
		errMsgs := utils.ValidateEmployeeExcel(request)
		if len(errMsgs) > 0 {
			detail := model.CreateEmployeeFailedDetail(i+1, errMsgs, &request)
			failed = append(failed, detail)
			continue
		}

		response, err := s.Create(ctx, userID, request)
		if err != nil {
			detail := model.CreateEmployeeFailedDetail(i+1, []string{err.Error()}, &request)
			failed = append(failed, detail)
			continue
		}

		success = append(success, *response)
		time.Sleep(100 * time.Millisecond)
	}

	response := model.CreateEmployeeBulkResponse(success, failed)
	return &response, nil
}

func (s *employeeService) GetReport(payload model.EmployeeBulkResponse) (*bytes.Buffer, error) {
	var success [][]interface{}
	for _, data := range payload.SuccessData {
		var secondaryLat, secondaryLng interface{}
		if data.CoordSecondary != nil {
			secondaryLat = data.CoordSecondary.Lat
			secondaryLng = data.CoordSecondary.Lng
		}

		success = append(success, []interface{}{
			data.EmployeeID, data.Name, data.Position, data.Department, data.Email, data.Phone,
			data.AddressPrimary, data.CoordPrimary.Lat, data.CoordPrimary.Lng,
			data.AddressSecondary, secondaryLat, secondaryLng,
			"SUCCESS", "",
		})
	}

	var failed [][]interface{}
	for _, data := range payload.FailedData {
		var primaryLat, primaryLng, secondaryLat, secondaryLng interface{}
		if data.Value.CoordPrimary != nil {
			primaryLat = data.Value.CoordPrimary.Lat
			primaryLng = data.Value.CoordPrimary.Lng
		}
		if data.Value.CoordSecondary != nil {
			secondaryLat = data.Value.CoordSecondary.Lat
			secondaryLng = data.Value.CoordSecondary.Lng
		}

		errMsg := strings.Join(data.Error, ", ")
		failed = append(failed, []interface{}{
			data.Value.EmployeeID, data.Value.Name, data.Value.Position, data.Value.Department, data.Value.Email, data.Value.Phone,
			data.Value.AddressPrimary, primaryLat, primaryLng,
			data.Value.AddressSecondary, secondaryLat, secondaryLng,
			"FAILED", errMsg,
		})
	}

	columns := utils.EmployeeReportColumns
	excel, err := utils.GenerateReport(columns, success, failed)
	if err != nil { return nil, apperror.ErrInternal("Failed to generate Excel") }

	return excel, nil
}

func (s *employeeService) GetTemplate() (*bytes.Buffer, error) {
	columns := utils.EmployeeColumns
	excel, err := utils.GenerateTemplate(columns)
	if err != nil { return nil, apperror.ErrInternal("Failed to generate Excel") }

	return excel, nil
}

func (s *employeeService) Export(ctx context.Context) (*bytes.Buffer, error) {
	employees, err := s.repository.FindAll(ctx)
	if err != nil { return nil, apperror.ErrInternal("Failed to find employees") }

	var data [][]interface{}
	for _, employee := range employees {
		var secondaryLat, secondaryLng interface{}
		if employee.CoordSecondary != nil {
			secondaryLat = employee.CoordSecondary.Lat
			secondaryLng = employee.CoordSecondary.Lng
		}

		data = append(data, []interface{}{
			employee.EmployeeID, employee.Name, employee.Position, employee.Department, employee.Email, employee.Phone,
			employee.AddressPrimary, employee.CoordPrimary.Lat, employee.CoordPrimary.Lng,
			employee.AddressSecondary, secondaryLat, secondaryLng,
		})
	}

	columns := utils.EmployeeColumns
	excel, err := utils.GenerateExport(columns, data)
	if err != nil { return nil, apperror.ErrInternal("Failed to generate Excel") }

	return excel, nil
}