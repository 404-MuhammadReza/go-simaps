package model

import "time"

// Hospital Model
type Hospital struct {
	ID         string     `bson:"_id"`
	Name       string     `bson:"name"`
	Phone	   string     `bson:"phone,omitempty"`
	Address    string     `bson:"address"`
	Coordinate Coordinate `bson:"coordinate"`
	CreatedAt  time.Time  `bson:"created_at"`
	UpdatedAt  time.Time  `bson:"updated_at"`
}

// DTOs
// Hospital Request DTO
type HospitalRequest struct {
	Name       string      `json:"name" binding:"required"`
	Phone      string      `json:"phone"`
	Address    string      `json:"address" binding:"required"`
	Coordinate *Coordinate `json:"coordinate"`
}

// Hospital Response DTO
type HospitalResponse struct {
	ID         string     `json:"id"`
	Name       string     `json:"name"`
	Phone      string     `json:"phone"`
	Address    string     `json:"address"`
	Coordinate Coordinate `json:"coordinate"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}

// Mapping functions
// Request to Model Hospital
func (request *HospitalRequest) ToModel(hospitalID string, coordinate Coordinate) Hospital {
	return Hospital{
		ID: hospitalID,
		Name: request.Name,
		Phone: request.Phone,
		Address: request.Address,
		Coordinate: coordinate,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// Model Hospital to Response
func (hospital *Hospital) ToResponse() HospitalResponse {
	return HospitalResponse{
		ID: hospital.ID,
		Name: hospital.Name,
		Phone: hospital.Phone,
		Address: hospital.Address,
		Coordinate: hospital.Coordinate,
		CreatedAt: hospital.CreatedAt,
		UpdatedAt: hospital.UpdatedAt,
	}
}

// Put update for Hospital model
func (hospital *Hospital) PutFrom(request HospitalRequest, coordinate Coordinate) {
	hospital.Name = request.Name
	hospital.Phone = request.Phone
	hospital.Address = request.Address
	hospital.Coordinate = coordinate
	hospital.UpdatedAt = time.Now()
}

// DTOs and Factory Functions for Bulk Operations
// Hospital Bulk Response
type HospitalBulkResponse struct {
	Total        int                    `json:"total_processed"`
	SuccessCount int                    `json:"success_count"`
	FailedCount  int                    `json:"failed_count"`
	SuccessData  []HospitalResponse     `json:"success"`
	FailedData   []HospitalFailedDetail `json:"failed"`
}

// Hospital Failed Detail
type HospitalFailedDetail struct {
	ROW   int              `json:"row"`
	Error []string         `json:"error"`
	Value *HospitalRequest `json:"value"`
}

// Create Hospital Bulk Response
func CreateHospitalBulkResponse(success []HospitalResponse, failed []HospitalFailedDetail) HospitalBulkResponse {
	return HospitalBulkResponse{
		Total: len(success) + len(failed),
		SuccessCount: len(success),
		FailedCount: len(failed),
		SuccessData: success,
		FailedData: failed,
	}
}

// Create Hospital Failed Detail
func CreateHospitalFailedDetail(row int, errMsgs []string, value *HospitalRequest) HospitalFailedDetail {
	return HospitalFailedDetail{
		ROW: row,
		Error: errMsgs,
		Value: value,
	}
}
