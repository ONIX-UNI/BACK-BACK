package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"mime"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"github.com/DuvanRozoParra/sicou/internal/appointments/repository"
)

func (s *Service) saveAttachments(
	ctx context.Context,
	identityDocument *multipart.FileHeader,
	utilityBill *multipart.FileHeader,
) ([]repository.AttachmentRecord, error) {
	candidates := []struct {
		fieldName string
		file      *multipart.FileHeader
	}{
		{fieldName: "identityDocument", file: identityDocument},
		{fieldName: "utilityBill", file: utilityBill},
	}

	hasFiles := false
	for _, candidate := range candidates {
		if candidate.file != nil {
			hasFiles = true
			break
		}
	}
	if !hasFiles {
		return nil, nil
	}

	if s.fileObjects == nil {
		return nil, errors.New("documents file service is not initialized")
	}

	baseURL := attachmentBaseURL()

	attachments := make([]repository.AttachmentRecord, 0, len(candidates))
	for _, candidate := range candidates {
		if candidate.file == nil {
			continue
		}

		file, err := candidate.file.Open()
		if err != nil {
			s.cleanupAttachments(ctx, attachments)
			return nil, err
		}

		mimeType := strings.TrimSpace(candidate.file.Header.Get("Content-Type"))
		if mimeType == "" {
			mimeType = mime.TypeByExtension(filepath.Ext(candidate.file.Filename))
		}
		if mimeType == "" {
			mimeType = "application/octet-stream"
		}

		fileObj, err := s.fileObjects.Create(
			ctx,
			file,
			candidate.file.Size,
			candidate.file.Filename,
			mimeType,
		)
		_ = file.Close()
		if err != nil {
			s.cleanupAttachments(ctx, attachments)
			return nil, err
		}

		if fileObj == nil || strings.TrimSpace(fileObj.ID) == "" {
			s.cleanupAttachments(ctx, attachments)
			return nil, errors.New("documents file service returned an empty file object id")
		}
		if strings.TrimSpace(fileObj.StorageKey) == "" {
			s.cleanupAttachments(ctx, attachments)
			return nil, errors.New("documents file service returned an empty storage key")
		}

		fileURL := fmt.Sprintf("%s/%s/download", baseURL, fileObj.ID)

		attachments = append(attachments, repository.AttachmentRecord{
			FileObjectID: fileObj.ID,
			FieldName:    candidate.fieldName,
			OriginalName: candidate.file.Filename,
			MimeType:     mimeType,
			SizeBytes:    candidate.file.Size,
			URL:          fileURL,
			StorageKey:   strings.TrimSpace(fileObj.StorageKey),
		})
	}

	return attachments, nil
}

func (s *Service) cleanupAttachments(ctx context.Context, attachments []repository.AttachmentRecord) {
	if s.fileObjects == nil {
		return
	}

	for _, attachment := range attachments {
		fileObjectID := strings.TrimSpace(attachment.FileObjectID)
		if fileObjectID == "" {
			continue
		}
		_ = s.fileObjects.Delete(ctx, fileObjectID)
	}
}

func attachmentBaseURL() string {
	baseURL := strings.TrimSpace(os.Getenv("FORM_LEGAL_ADVISE_PUBLIC_BASE_URL"))
	if baseURL == "" {
		baseURL = strings.TrimSpace(os.Getenv("ASESORIAS_PUBLIC_BASE_URL"))
	}
	if baseURL == "" {
		baseURL = "/api/v1/documents/files"
	}
	return strings.TrimRight(strings.ReplaceAll(baseURL, "\\", "/"), "/")
}

func buildPayload(
	in CreateAsesoriaInput,
	attachments []repository.AttachmentRecord,
) (string, error) {
	attachmentsPayload := make([]map[string]any, 0, len(attachments))
	for _, attachment := range attachments {
		attachmentsPayload = append(attachmentsPayload, map[string]any{
			"fieldName": attachment.FieldName,
			"name":      attachment.OriginalName,
			"type":      attachment.MimeType,
			"size":      attachment.SizeBytes,
			"url":       attachment.URL,
		})
	}

	payload := map[string]any{
		"acceptsDataProcessing":  in.AcceptsDataProcessing,
		"consultationDate":       in.ConsultationDate,
		"fullName":               in.FullName,
		"documentType":           in.DocumentType,
		"otherDocumentType":      in.OtherDocumentType,
		"documentNumber":         in.DocumentNumber,
		"birthDate":              in.BirthDate,
		"age":                    in.Age,
		"maritalStatus":          in.MaritalStatus,
		"otherMaritalStatus":     in.OtherMaritalStatus,
		"gender":                 in.Gender,
		"address":                in.Address,
		"housingType":            in.HousingType,
		"otherHousingType":       in.OtherHousingType,
		"socioEconomicStratum":   in.SocioEconomicStratum,
		"sisbenCategory":         in.SisbenCategory,
		"mobilePhone":            in.MobilePhone,
		"email":                  in.Email,
		"populationType":         in.PopulationType,
		"otherPopulationType":    in.OtherPopulationType,
		"headOfHousehold":        in.HeadOfHousehold,
		"occupation":             in.Occupation,
		"educationLevel":         in.EducationLevel,
		"otherEducationLevel":    in.OtherEducationLevel,
		"caseDescription":        in.CaseDescription,
		"authorizesNotification": in.AuthorizesNotification,
		"submittedAt":            in.SubmittedAt,
		"attachments":            attachmentsPayload,
	}

	bytes, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}
