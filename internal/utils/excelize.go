package utils

import (
	"io"
	"bytes"
	"errors"
	"strconv"
	"strings"

	"go-simaps/internal/model"

	"github.com/xuri/excelize/v2"
)

type Column struct {
	Header string
	Width  float64
	Align  string
}

var (
	EmployeeColumns = []Column{
		{Header: "Employee ID", Width: 15, Align: "center"},
		{Header: "Full Name", Width: 30, Align: "left"},
		{Header: "Position", Width: 30, Align: "left"},
		{Header: "Department", Width: 30, Align: "left"},
		{Header: "Email", Width: 30, Align: "left"},
		{Header: "Phone Number", Width: 15, Align: "center"},
		{Header: "KTP Address", Width: 80, Align: "left"},
		{Header: "KTP Latitude", Width: 20, Align: "center"},
		{Header: "KTP Longitude", Width: 20, Align: "center"},
		{Header: "Domicile Address", Width: 80, Align: "left"},
		{Header: "Domicile Latitude", Width: 20, Align: "center"},
		{Header: "Domicile Longitude", Width: 20, Align: "center"},
	}

	EmployeeReportColumns = []Column{
		{Header: "Employee ID", Width: 15, Align: "center"},
		{Header: "Full Name", Width: 30, Align: "left"},
		{Header: "Position", Width: 30, Align: "left"},
		{Header: "Department", Width: 30, Align: "left"},
		{Header: "Email", Width: 30, Align: "left"},
		{Header: "Phone Number", Width: 15, Align: "center"},
		{Header: "KTP Address", Width: 80, Align: "left"},
		{Header: "KTP Latitude", Width: 20, Align: "center"},
		{Header: "KTP Longitude", Width: 20, Align: "center"},
		{Header: "Domicile Address", Width: 80, Align: "left"},
		{Header: "Domicile Latitude", Width: 20, Align: "center"},
		{Header: "Domicile Longitude", Width: 20, Align: "center"},
		{Header: "Status", Width: 15, Align: "center"},
		{Header: "Description", Width: 50, Align: "left"},
	}

	HospitalColumns = []Column{
		{Header: "Name", Width: 30, Align: "left"},
		{Header: "Phone Number", Width: 15, Align: "center"},
		{Header: "Address", Width: 80, Align: "left"},
		{Header: "Latitude", Width: 15, Align: "center"},
		{Header: "Longitude", Width: 15, Align: "center"},
	}

	HospitalReportColumns = []Column{
		{Header: "Name", Width: 30, Align: "left"},
		{Header: "Phone Number", Width: 15, Align: "center"},
		{Header: "Address", Width: 80, Align: "left"},
		{Header: "Latitude", Width: 15, Align: "center"},
		{Header: "Longitude", Width: 15, Align: "center"},
		{Header: "Status", Width: 15, Align: "center"},
		{Header: "Description", Width: 50, Align: "left"},
	}
)

var (
	border = []excelize.Border{
		{Type: "left", Color: "000000", Style: 1},
		{Type: "top", Color: "000000", Style: 1},
		{Type: "bottom", Color: "000000", Style: 1},
		{Type: "right", Color: "000000", Style: 1},
	}

	alignCenter = &excelize.Alignment{
		Horizontal: "center",
		Vertical:   "center",
		WrapText:   true,
	}

	alignLeft = &excelize.Alignment{
		Horizontal: "left",
		Vertical:   "center",
		WrapText:   true,
	}

	header = &excelize.Style{
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"#ededed"},
			Pattern: 1,
		},
		Font:      &excelize.Font{Bold: true},
		Border:    border,
		Alignment: alignCenter,
	}

	textCenter = &excelize.Style{
		NumFmt:    49,
		Border:    border,
		Alignment: alignCenter,
	}

	textLeft = &excelize.Style{
		NumFmt:    49,
		Border:    border,
		Alignment: alignLeft,
	}

	successCenter = &excelize.Style{
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"#C6EFCE"},
			Pattern: 1,
		},
		Font:      &excelize.Font{Color: "#006100"},
		Border:    border,
		Alignment: alignCenter,
	}

	successLeft = &excelize.Style{
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"#C6EFCE"},
			Pattern: 1,
		},
		Font:      &excelize.Font{Color: "#006100"},
		Border:    border,
		Alignment: alignLeft,
	}

	failedCenter = &excelize.Style{
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"#FFC7CE"},
			Pattern: 1,
		},
		Font:      &excelize.Font{Color: "#9C0006"},
		Border:    border,
		Alignment: alignCenter,
	}

	failedLeft = &excelize.Style{
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"#FFC7CE"},
			Pattern: 1,
		},
		Font:      &excelize.Font{Color: "#9C0006"},
		Border:    border,
		Alignment: alignLeft,
	}
)

func ReadExcel(file io.Reader) ([][]string, error) {
	excel, err := excelize.OpenReader(file)
	if err != nil { return nil, err }
	defer excel.Close()

	rows, err := excel.GetRows("Sheet1")
	if err != nil { return nil, err }

	return rows, nil
}

func ParseEmployeeExcel(file io.Reader) ([]model.EmployeeRequest, error) {
	rows, err := ReadExcel(file)
	if err != nil { return nil, err }
	if len(rows) == 0 { return nil, errors.New("Excel file is empty") }
	if len(rows[0]) != len(EmployeeColumns) { return nil, errors.New("Excel columns do not match expected format") }

	var employees []model.EmployeeRequest
	for i, row := range rows {
		if i == 0 { continue }

		var request model.EmployeeRequest
		var primaryLat, primaryLng, secondaryLat, secondaryLng string

		if len(row) > 0 { request.EmployeeID = strings.TrimSpace(row[0]) }
		if len(row) > 1 { request.Name = strings.TrimSpace(row[1]) }
		if len(row) > 2 { request.Position = strings.TrimSpace(row[2]) }
		if len(row) > 3 { request.Department = strings.TrimSpace(row[3]) }
		if len(row) > 4 { request.Email = strings.TrimSpace(row[4]) }
		if len(row) > 5 { request.Phone = strings.TrimSpace(row[5]) }

		if len(row) > 6 { request.AddressPrimary = strings.TrimSpace(row[6]) }
		if len(row) > 7 { primaryLat = strings.TrimSpace(row[7]) }
		if len(row) > 8 { primaryLng = strings.TrimSpace(row[8]) }

		if primaryLat != "" && primaryLng != "" {
			latVal, err1 := strconv.ParseFloat(strings.ReplaceAll(primaryLat, ",", "."), 64)
			lngVal, err2 := strconv.ParseFloat(strings.ReplaceAll(primaryLng, ",", "."), 64)
			if err1 == nil && err2 == nil { request.CoordPrimary = &model.Coordinate{Lat: latVal, Lng: lngVal} }
		}

		if len(row) > 9 { request.AddressSecondary = strings.TrimSpace(row[9]) }
		if len(row) > 10 { secondaryLat = strings.TrimSpace(row[10]) }
		if len(row) > 11 { secondaryLng = strings.TrimSpace(row[11]) }

		if secondaryLat != "" && secondaryLng != "" {
			latVal, err1 := strconv.ParseFloat(strings.ReplaceAll(secondaryLat, ",", "."), 64)
			lngVal, err2 := strconv.ParseFloat(strings.ReplaceAll(secondaryLng, ",", "."), 64)
			if err1 == nil && err2 == nil { request.CoordSecondary = &model.Coordinate{Lat: latVal, Lng: lngVal} }
		}

		employees = append(employees, request)
	}

	return employees, nil
}

func ParseHospitalExcel(file io.Reader) ([]model.HospitalRequest, error) {
	rows, err := ReadExcel(file)
	if err != nil { return nil, err }
	if len(rows) == 0 { return nil, errors.New("Excel file is empty") }
	if len(rows[0]) != len(HospitalColumns) { return nil, errors.New("Excel columns do not match expected format") }

	var hospitals []model.HospitalRequest
	for i, row := range rows {
		if i == 0 { continue }

		var lat, lng string
		var request model.HospitalRequest
		if len(row) > 0 { request.Name = strings.TrimSpace(row[0]) }
		if len(row) > 1 { request.Phone = strings.TrimSpace(row[1]) }
		if len(row) > 2 { request.Address = strings.TrimSpace(row[2]) }
		if len(row) > 3 { lat = strings.TrimSpace(row[3]) }
		if len(row) > 4 { lng = strings.TrimSpace(row[4]) }

		if lat != "" && lng != "" {
			latVal, err1 := strconv.ParseFloat(strings.ReplaceAll(lat, ",", "."), 64)
			lngVal, err2 := strconv.ParseFloat(strings.ReplaceAll(lng, ",", "."), 64)
			if err1 == nil && err2 == nil { request.Coordinate = &model.Coordinate{Lat: latVal, Lng: lngVal} }
		}

		hospitals = append(hospitals, request)
	}

	return hospitals, nil
}

func ValidateEmployeeExcel(request model.EmployeeRequest) []string {
	var errMsgs []string
	if request.Name == "" { errMsgs = append(errMsgs, "Name is required") }
	if request.Position == "" { errMsgs = append(errMsgs, "Position is required") }
	if request.Phone == "" { errMsgs = append(errMsgs, "Phone number is required") }
	if request.EmployeeID == "" { errMsgs = append(errMsgs, "Employee ID is required") }
	if request.Email == "" { errMsgs = append(errMsgs, "Email is required") }
	if request.AddressPrimary == "" { errMsgs = append(errMsgs, "KTP Address is required") }

	return errMsgs
}

func ValidateHospitalExcel(request model.HospitalRequest) []string {
	var errMsgs []string
	if request.Name == "" { errMsgs = append(errMsgs, "Name is required") }
	if request.Address == "" { errMsgs = append(errMsgs, "Address is required") }

	return errMsgs
}

func GenerateTemplate(columns []Column) (*bytes.Buffer, error) {
	excel := excelize.NewFile()
	defer excel.Close()

	sheet := "Sheet1"
	styleHeader, _ := excel.NewStyle(header)
	styleTextCenter, _ := excel.NewStyle(textCenter)
	styleTextLeft, _ := excel.NewStyle(textLeft)

	for i, col := range columns {
		colName, _ := excelize.ColumnNumberToName(i + 1)
		cellName := colName + "1"

		excel.SetCellValue(sheet, cellName, col.Header)
		excel.SetColWidth(sheet, colName, colName, col.Width)

		colRange := colName + ":" + colName
		switch col.Align {
		case "center": excel.SetColStyle(sheet, colRange, styleTextCenter)
		case "left": excel.SetColStyle(sheet, colRange, styleTextLeft)
		}
	}

	lastColName, _ := excelize.ColumnNumberToName(len(columns))
	excel.SetCellStyle(sheet, "A1", lastColName+"1", styleHeader)

	var buf bytes.Buffer
	if err := excel.Write(&buf); err != nil { return nil, err }

	return &buf, nil
}

func GenerateReport(columns []Column, success, failed [][]interface{}) (*bytes.Buffer, error) {
	excel := excelize.NewFile()
	defer excel.Close()

	sheet := "Sheet1"
	styleHeader, _ := excel.NewStyle(header)
	styleSuccessLeft, _ := excel.NewStyle(successLeft)
	styleSuccessCenter, _ := excel.NewStyle(successCenter)
	styleFailedLeft, _ := excel.NewStyle(failedLeft)
	styleFailedCenter, _ := excel.NewStyle(failedCenter)

	for i, col := range columns {
		cellName, _ := excelize.CoordinatesToCellName(i+1, 1)
		excel.SetCellValue(sheet, cellName, col.Header)

		colName, _ := excelize.ColumnNumberToName(i + 1)
		excel.SetColWidth(sheet, colName, colName, col.Width)
	}

	lastColName, _ := excelize.ColumnNumberToName(len(columns))
	excel.SetCellStyle(sheet, "A1", lastColName+"1", styleHeader)

	rowNum := 2
	writeRows := func(data [][]interface{}, styleLeft, styleCenter int) {
		for _, row := range data {
			cellStart, _ := excelize.CoordinatesToCellName(1, rowNum)
			excel.SetSheetRow(sheet, cellStart, &row)

			rowStr := strconv.Itoa(rowNum)
			for i := range columns {
				colName, _ := excelize.ColumnNumberToName(i + 1)
				targetCell := colName + rowStr
				switch columns[i].Align {
				case "center": excel.SetCellStyle(sheet, targetCell, targetCell, styleCenter)
				case "left": excel.SetCellStyle(sheet, targetCell, targetCell, styleLeft)
				}
			}
			rowNum++
		}
	}

	writeRows(success, styleSuccessLeft, styleSuccessCenter)
	writeRows(failed, styleFailedLeft, styleFailedCenter)

	var buf bytes.Buffer
	if err := excel.Write(&buf); err != nil { return nil, err }

	return &buf, nil
}

func GenerateExport(columns []Column, data [][]interface{}) (*bytes.Buffer, error) {
	excel := excelize.NewFile()
	defer excel.Close()

	sheet := "Sheet1"
	styleHeader, _ := excel.NewStyle(header)
	styleTextCenter, _ := excel.NewStyle(textCenter)
	styleTextLeft, _ := excel.NewStyle(textLeft)

	for i, col := range columns {
		colName, _ := excelize.ColumnNumberToName(i + 1)
		celName := colName + "1"

		excel.SetCellValue(sheet, celName, col.Header)
		excel.SetColWidth(sheet, colName, colName, col.Width)
	}

	row := 2
	for _, record := range data {
		cellStart, _ := excelize.CoordinatesToCellName(1, row)
		excel.SetSheetRow(sheet, cellStart, &record)

		rowStr := strconv.Itoa(row)
		for i := range columns {
			colName, _ := excelize.ColumnNumberToName(i + 1)
			targetCell := colName + rowStr
			switch columns[i].Align {
			case "center": excel.SetCellStyle(sheet, targetCell, targetCell, styleTextCenter)
			case "left": excel.SetCellStyle(sheet, targetCell, targetCell, styleTextLeft)
			}
		}
		row++
	}

	lastColName, _ := excelize.ColumnNumberToName(len(columns))
	excel.SetCellStyle(sheet, "A1", lastColName+"1", styleHeader)

	var buf bytes.Buffer
	if err := excel.Write(&buf); err != nil { return nil, err }

	return &buf, nil
}
