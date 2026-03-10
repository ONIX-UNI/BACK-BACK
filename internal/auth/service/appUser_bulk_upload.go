package service

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"sort"
	"strconv"
	"strings"

	"github.com/DuvanRozoParra/sicou/internal/auth/dto"
	"github.com/google/uuid"
)

var requiredBulkHeaders = []string{
	"display_name",
	"email",
	"password_hash",
	"role",
	"is_active",
}

const (
	maxBulkXLSXBytes int64 = 10 * 1024 * 1024
	maxZipEntryBytes int64 = 8 * 1024 * 1024
)

func (s *AppUserService) BulkCreate(ctx context.Context, file io.Reader, actorUserID *uuid.UUID) (*dto.BulkCreateAppUserResponse, error) {
	rows, err := parseXLSXRows(file)
	if err != nil {
		return nil, fmt.Errorf("%w: archivo xlsx invalido", ErrInvalidInput)
	}
	if len(rows) == 0 {
		return nil, fmt.Errorf("%w: el archivo esta vacio", ErrInvalidInput)
	}

	headerIndex, err := mapBulkHeaderIndexes(rows[0])
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrInvalidInput, err.Error())
	}

	response := &dto.BulkCreateAppUserResponse{
		Summary: dto.BulkCreateSummary{},
		Results: make([]dto.BulkCreateResultItem, 0, maxInt(0, len(rows)-1)),
	}

	seenEmails := make(map[string]int)
	for i := 1; i < len(rows); i++ {
		line := i + 1
		row := rows[i]
		if isEmptyRow(row) {
			continue
		}
		response.Summary.Total++

		item, created := s.bulkCreateRow(ctx, row, line, headerIndex, seenEmails, actorUserID)
		response.Results = append(response.Results, item)
		if created {
			response.Summary.Created++
		} else {
			response.Summary.Failed++
		}
	}

	return response, nil
}

func (s *AppUserService) bulkCreateRow(
	ctx context.Context,
	row []string,
	line int,
	headerIndex map[string]int,
	seenEmails map[string]int,
	actorUserID *uuid.UUID,
) (dto.BulkCreateResultItem, bool) {
	req := dto.CreateAppUserRequest{
		DisplayName: strings.TrimSpace(cellByHeader(row, headerIndex, "display_name")),
		Email:       strings.TrimSpace(cellByHeader(row, headerIndex, "email")),
		PasswordHash: strings.TrimSpace(
			cellByHeader(row, headerIndex, "password_hash"),
		),
		Role: strings.TrimSpace(cellByHeader(row, headerIndex, "role")),
	}

	isActiveRaw := strings.TrimSpace(cellByHeader(row, headerIndex, "is_active"))
	if isActiveRaw != "" {
		parsed, err := parseBoolCell(isActiveRaw)
		if err != nil {
			return dto.BulkCreateResultItem{
				Line:    line,
				Email:   req.Email,
				Status:  "error",
				Message: "Valor invalido en is_active",
			}, false
		}
		req.IsActive = &parsed
	}

	if req.DisplayName == "" || req.Email == "" || req.PasswordHash == "" || req.Role == "" {
		return dto.BulkCreateResultItem{
			Line:    line,
			Email:   req.Email,
			Status:  "error",
			Message: "Campos obligatorios incompletos",
		}, false
	}

	normalizedEmail := strings.ToLower(req.Email)
	if previousLine, exists := seenEmails[normalizedEmail]; exists {
		return dto.BulkCreateResultItem{
			Line:    line,
			Email:   req.Email,
			Status:  "error",
			Message: "Correo duplicado en archivo (linea " + strconv.Itoa(previousLine) + ")",
		}, false
	}
	seenEmails[normalizedEmail] = line

	existing, err := s.repo.GetByEmail(ctx, req.Email)
	if err != nil {
		return dto.BulkCreateResultItem{
			Line:    line,
			Email:   req.Email,
			Status:  "error",
			Message: "Error consultando correo en base de datos",
		}, false
	}
	if existing != nil {
		return dto.BulkCreateResultItem{
			Line:    line,
			Email:   req.Email,
			Status:  "error",
			Message: "Correo ya existe",
		}, false
	}

	user, err := s.Create(ctx, req, actorUserID)
	if err != nil {
		message := "No se pudo crear el usuario"
		switch {
		case err == ErrRoleNotFound:
			message = "Rol no existe"
		case err == ErrEmailAlreadyUsed:
			message = "Correo ya existe"
		case err == ErrInvalidInput:
			message = "Datos invalidos"
		case strings.Contains(strings.ToLower(err.Error()), "duplicate key"),
			strings.Contains(strings.ToLower(err.Error()), "unique"):
			message = "Correo ya existe"
		}
		return dto.BulkCreateResultItem{
			Line:    line,
			Email:   req.Email,
			Status:  "error",
			Message: message,
		}, false
	}

	return dto.BulkCreateResultItem{
		Line:   line,
		Email:  user.Email,
		Status: "created",
	}, true
}

func mapBulkHeaderIndexes(headers []string) (map[string]int, error) {
	indexes := make(map[string]int, len(headers))
	for i, header := range headers {
		normalized := strings.ToLower(strings.TrimSpace(header))
		if normalized == "" {
			continue
		}
		indexes[normalized] = i
	}

	missing := make([]string, 0, len(requiredBulkHeaders))
	for _, required := range requiredBulkHeaders {
		if _, ok := indexes[required]; !ok {
			missing = append(missing, required)
		}
	}
	if len(missing) > 0 {
		sort.Strings(missing)
		return nil, fmt.Errorf("faltan columnas requeridas: %s", strings.Join(missing, ", "))
	}

	return indexes, nil
}

func parseBoolCell(value string) (bool, error) {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "1", "true", "t", "si", "s", "yes", "y":
		return true, nil
	case "0", "false", "f", "no", "n":
		return false, nil
	default:
		return false, ErrInvalidInput
	}
}

func cellByHeader(row []string, headerIndex map[string]int, header string) string {
	idx, ok := headerIndex[header]
	if !ok || idx < 0 || idx >= len(row) {
		return ""
	}
	return row[idx]
}

func parseXLSXRows(file io.Reader) ([][]string, error) {
	limited := &io.LimitedReader{R: file, N: maxBulkXLSXBytes + 1}
	content, err := io.ReadAll(limited)
	if err != nil {
		return nil, err
	}
	if int64(len(content)) > maxBulkXLSXBytes {
		return nil, ErrFileTooLarge
	}
	if len(content) == 0 {
		return nil, ErrInvalidInput
	}

	readerAt := bytes.NewReader(content)
	zipReader, err := zip.NewReader(readerAt, int64(len(content)))
	if err != nil {
		return nil, err
	}

	files := make(map[string]*zip.File, len(zipReader.File))
	worksheetNames := make([]string, 0)
	for _, zf := range zipReader.File {
		files[zf.Name] = zf
		if strings.HasPrefix(zf.Name, "xl/worksheets/") && strings.HasSuffix(zf.Name, ".xml") {
			worksheetNames = append(worksheetNames, zf.Name)
		}
	}
	if len(worksheetNames) == 0 {
		return nil, ErrInvalidInput
	}
	sort.Strings(worksheetNames)

	sharedStrings := make([]string, 0)
	if sharedFile, ok := files["xl/sharedStrings.xml"]; ok {
		sharedStrings, err = parseSharedStrings(sharedFile)
		if err != nil {
			return nil, err
		}
	}

	rows, err := parseWorksheet(files[worksheetNames[0]], sharedStrings)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

type xlsxSST struct {
	SI []xlsxSI `xml:"si"`
}

type xlsxSI struct {
	T string     `xml:"t"`
	R []xlsxText `xml:"r"`
}

type xlsxText struct {
	T string `xml:"t"`
}

func parseSharedStrings(zf *zip.File) ([]string, error) {
	content, err := readZipFile(zf, maxZipEntryBytes)
	if err != nil {
		return nil, err
	}

	var sst xlsxSST
	if err := xml.Unmarshal(content, &sst); err != nil {
		return nil, err
	}

	values := make([]string, 0, len(sst.SI))
	for _, si := range sst.SI {
		if strings.TrimSpace(si.T) != "" || len(si.R) == 0 {
			values = append(values, si.T)
			continue
		}

		var builder strings.Builder
		for _, fragment := range si.R {
			builder.WriteString(fragment.T)
		}
		values = append(values, builder.String())
	}

	return values, nil
}

type xlsxWorksheet struct {
	SheetData xlsxSheetData `xml:"sheetData"`
}

type xlsxSheetData struct {
	Rows []xlsxRow `xml:"row"`
}

type xlsxRow struct {
	R int        `xml:"r,attr"`
	C []xlsxCell `xml:"c"`
}

type xlsxCell struct {
	R  string        `xml:"r,attr"`
	T  string        `xml:"t,attr"`
	V  string        `xml:"v"`
	IS *xlsxInlineSI `xml:"is"`
}

type xlsxInlineSI struct {
	T string     `xml:"t"`
	R []xlsxText `xml:"r"`
}

func parseWorksheet(zf *zip.File, sharedStrings []string) ([][]string, error) {
	content, err := readZipFile(zf, maxZipEntryBytes)
	if err != nil {
		return nil, err
	}

	var ws xlsxWorksheet
	if err := xml.Unmarshal(content, &ws); err != nil {
		return nil, err
	}

	maxRow := 0
	rowsByIndex := make(map[int][]string, len(ws.SheetData.Rows))
	for i, row := range ws.SheetData.Rows {
		rowNumber := row.R
		if rowNumber <= 0 {
			rowNumber = i + 1
		}
		if rowNumber > maxRow {
			maxRow = rowNumber
		}

		maxCol := -1
		values := make(map[int]string, len(row.C))
		for j, cell := range row.C {
			col := columnIndexFromCellRef(cell.R)
			if col < 0 {
				col = j
			}
			if col > maxCol {
				maxCol = col
			}

			values[col] = resolveCellValue(cell, sharedStrings)
		}
		if maxCol < 0 {
			rowsByIndex[rowNumber] = []string{}
			continue
		}

		rowValues := make([]string, maxCol+1)
		for col, value := range values {
			rowValues[col] = value
		}
		rowsByIndex[rowNumber] = rowValues
	}

	rows := make([][]string, maxRow)
	for rowNumber := 1; rowNumber <= maxRow; rowNumber++ {
		rows[rowNumber-1] = rowsByIndex[rowNumber]
	}

	return rows, nil
}

func resolveCellValue(cell xlsxCell, sharedStrings []string) string {
	switch strings.ToLower(strings.TrimSpace(cell.T)) {
	case "s":
		idx, err := strconv.Atoi(strings.TrimSpace(cell.V))
		if err != nil || idx < 0 || idx >= len(sharedStrings) {
			return ""
		}
		return sharedStrings[idx]
	case "inlineStr":
		if cell.IS == nil {
			return ""
		}
		if len(cell.IS.R) == 0 {
			return cell.IS.T
		}

		var builder strings.Builder
		for _, fragment := range cell.IS.R {
			builder.WriteString(fragment.T)
		}
		return builder.String()
	case "b":
		raw := strings.TrimSpace(cell.V)
		if raw == "1" {
			return "true"
		}
		if raw == "0" {
			return "false"
		}
		return raw
	default:
		return cell.V
	}
}

func columnIndexFromCellRef(cellRef string) int {
	ref := strings.TrimSpace(cellRef)
	if ref == "" {
		return -1
	}

	columnPart := make([]rune, 0, len(ref))
	for _, r := range ref {
		if r >= 'A' && r <= 'Z' {
			columnPart = append(columnPart, r)
			continue
		}
		if r >= 'a' && r <= 'z' {
			columnPart = append(columnPart, r-'a'+'A')
			continue
		}
		break
	}

	if len(columnPart) == 0 {
		return -1
	}

	index := 0
	for _, r := range columnPart {
		index = index*26 + int(r-'A'+1)
	}

	return index - 1
}

func readZipFile(zf *zip.File, maxBytes int64) ([]byte, error) {
	rc, err := zf.Open()
	if err != nil {
		return nil, err
	}
	defer rc.Close()

	limited := &io.LimitedReader{R: rc, N: maxBytes + 1}
	content, err := io.ReadAll(limited)
	if err != nil {
		return nil, err
	}
	if int64(len(content)) > maxBytes {
		return nil, ErrFileTooLarge
	}
	return content, nil
}

func maxInt(left int, right int) int {
	if left > right {
		return left
	}
	return right
}

func isEmptyRow(row []string) bool {
	for _, value := range row {
		if strings.TrimSpace(value) != "" {
			return false
		}
	}
	return true
}
