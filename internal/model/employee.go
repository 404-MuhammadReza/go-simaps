package model

import "time"

// Employee Model
type Employee struct {
	ID               string      `bson:"_id"`
	Name             string      `bson:"name"`
	Position         string      `bson:"title"`
	Department       string      `bson:"department,omitempty"`
	EmployeeID       string      `bson:"employee_id"`
	Email 		     string      `bson:"email"`

	Phone            string      `bson:"phone"`
	AddressPrimary   string      `bson:"address_primary"`
	CoordPrimary     Coordinate  `bson:"coordinate_primary"`
	AddressSecondary string      `bson:"address_secondary,omitempty"`
	CoordSecondary   *Coordinate `bson:"coordinate_secondary,omitempty"`

	Picture          string      `bson:"picture_object,omitempty"`
	CreatedAt        time.Time   `bson:"created_at"`
	UpdatedAt        time.Time   `bson:"updated_at"`
}

// Coordinate Sub-Model
type Coordinate struct {
	Lat float64 `bson:"lat" json:"lat" binding:"required"`
	Lng float64 `bson:"lng" json:"lng" binding:"required"`
}

// DTOs
// Employee Request DTO
type EmployeeRequest struct {
	Name             string      `json:"name" binding:"required"`
	Position         string      `json:"position" binding:"required"`
	Department       string      `json:"department"`
	EmployeeID       string      `json:"employee_id" binding:"required"`
	Email			 string      `json:"email" binding:"required,email"`

	Phone            string      `json:"phone" binding:"required"`
	AddressPrimary   string      `json:"address_primary" binding:"required"`
	CoordPrimary     *Coordinate `json:"coordinate_primary"`
	AddressSecondary string      `json:"address_secondary"`
	CoordSecondary   *Coordinate `json:"coordinate_secondary"`
}

// Employee Response DTO
type EmployeeResponse struct {
	ID               string      `json:"id"`
	Name             string      `json:"name"`
	Position         string      `json:"position"`
	Department       string      `json:"department"`
    EmployeeID       string      `json:"employee_id"`
	Email 		     string      `json:"email"`

	Phone            string      `json:"phone"`
	AddressPrimary   string      `json:"address_primary"`
	CoordPrimary     Coordinate  `json:"coordinate_primary"`
	AddressSecondary string      `json:"address_secondary"`
	CoordSecondary   *Coordinate `json:"coordinate_secondary"`

    PictureUrl       string      `json:"picture_url"`
	CreatedAt        time.Time   `json:"created_at"`
	UpdatedAt        time.Time   `json:"updated_at"`
}

// Mapping Functions
// Request to Model Employee
func (request *EmployeeRequest) ToModel(userID string, coordPrimary Coordinate, coordSecondary *Coordinate) Employee {
	return Employee{
		ID: userID,
		Name: request.Name,
		Position: request.Position,
		Department: request.Department,
		EmployeeID: request.EmployeeID,
		Email: request.Email,
		Phone: request.Phone,
		AddressPrimary: request.AddressPrimary,
		CoordPrimary: coordPrimary,
		AddressSecondary: request.AddressSecondary,
		CoordSecondary: coordSecondary,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// Model Employee to Response
func (employee *Employee) ToResponse(pictureUrl string) EmployeeResponse {
	return EmployeeResponse{
		ID: employee.ID,
		Name: employee.Name,
		Position: employee.Position,
		Department: employee.Department,
		Email: employee.Email,
		Phone: employee.Phone,
		EmployeeID: employee.EmployeeID,
		AddressPrimary: employee.AddressPrimary,
		CoordPrimary: employee.CoordPrimary,
		AddressSecondary: employee.AddressSecondary,
		CoordSecondary: employee.CoordSecondary,
		PictureUrl: pictureUrl,
		CreatedAt: employee.CreatedAt,
		UpdatedAt: employee.UpdatedAt,
	}
}

// Put update for Employee model
func (employee *Employee) PutFrom(request EmployeeRequest, coordPrimary Coordinate, coordSecondary *Coordinate) {
    employee.Name = request.Name
    employee.Position = request.Position
    employee.Department = request.Department
    employee.EmployeeID = request.EmployeeID
	employee.Email = request.Email

    employee.Phone = request.Phone
    employee.AddressPrimary = request.AddressPrimary
    employee.CoordPrimary = coordPrimary
    
    employee.AddressSecondary = request.AddressSecondary
	employee.CoordSecondary = coordSecondary

	employee.UpdatedAt = time.Now()
}

// DTOs and Factory Functions for Bulk Operations
// Employee Bulk Response
type EmployeeBulkResponse struct {
	Total        int                    `json:"total_processed"`
	SuccessCount int                    `json:"success_count"`
	FailedCount  int                    `json:"failed_count"`
	SuccessData  []EmployeeResponse     `json:"success"`
	FailedData   []EmployeeFailedDetail `json:"failed"`
}

// Employee Failed Detail
type EmployeeFailedDetail struct {
	ROW   int              `json:"row"`
	Error []string         `json:"error"`
	Value *EmployeeRequest `json:"value"`
}

// Create Employee Bulk Response
func CreateEmployeeBulkResponse(success []EmployeeResponse, failed []EmployeeFailedDetail) EmployeeBulkResponse {
	return EmployeeBulkResponse{
		Total: len(success) + len(failed),
		SuccessCount: len(success),
		FailedCount: len(failed),
		SuccessData: success,
		FailedData: failed,
	}
}

// Create Employee Failed Detail
func CreateEmployeeFailedDetail(row int, errMsgs []string, value *EmployeeRequest) EmployeeFailedDetail {
	return EmployeeFailedDetail{
		ROW: row,
		Error: errMsgs,
		Value: value,
	}
}
