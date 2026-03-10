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

	"github.com/DuvanRozoParra/sicou/internal/pqrs/repository"
)

func (s *Service) saveAttachments(ctx context.Context, files []*multipart.FileHeader) ([]repository.AttachmentRecord, error) {
	if len(files) == 0 {
		return nil, nil
	}

	if s.fileObjects == nil {
		return nil, errors.New("documents file service is not initialized")
	}

	baseURL := attachmentBaseURL()

	attachments := make([]repository.AttachmentRecord, 0, len(files))
	for _, fileHeader := range files {
		if fileHeader == nil {
			continue
		}

		file, err := fileHeader.Open()
		if err != nil {
			s.cleanupAttachments(ctx, attachments)
			return nil, err
		}

		mimeType := strings.TrimSpace(fileHeader.Header.Get("Content-Type"))
		if mimeType == "" {
			mimeType = mime.TypeByExtension(filepath.Ext(fileHeader.Filename))
		}
		if mimeType == "" {
			mimeType = "application/octet-stream"
		}

		fileObj, err := s.fileObjects.Create(
			ctx,
			file,
			fileHeader.Size,
			fileHeader.Filename,
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

		fileURL := fmt.Sprintf("%s/%s/download", baseURL, fileObj.ID)

		attachments = append(attachments, repository.AttachmentRecord{
			FileObjectID: fileObj.ID,
			OriginalName: fileHeader.Filename,
			MimeType:     mimeType,
			SizeBytes:    fileHeader.Size,
			URL:          fileURL,
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
	baseURL := strings.TrimSpace(os.Getenv("PQRS_PUBLIC_BASE_URL"))
	if baseURL == "" {
		baseURL = "/api/v1/documents/files"
	}
	return strings.TrimRight(strings.ReplaceAll(baseURL, "\\", "/"), "/")
}

func buildBodyPayload(
	in CreatePQRSInput,
	attachments []repository.AttachmentRecord,
) (string, error) {
	attachmentsPayload := make([]map[string]any, 0, len(attachments))
	for _, attachment := range attachments {
		attachmentsPayload = append(attachmentsPayload, map[string]any{
			"name": attachment.OriginalName,
			"type": attachment.MimeType,
			"size": attachment.SizeBytes,
			"url":  attachment.URL,
		})
	}

	payload := map[string]any{
		"queryType":                 in.QueryType,
		"personType":                in.PersonType,
		"documentType":              in.DocumentType,
		"documentOrTaxId":           in.DocumentOrTaxID,
		"firstName":                 in.FirstName,
		"middleName":                in.MiddleName,
		"firstLastName":             in.FirstLastName,
		"secondLastName":            in.SecondLastName,
		"gender":                    in.Gender,
		"address":                   in.Address,
		"neighborhoodArea":          in.NeighborhoodArea,
		"allowsElectronicResponse":  in.AllowsElectronicResponse,
		"email":                     in.Email,
		"phone":                     in.Phone,
		"populationGroup":           in.PopulationGroup,
		"otherPopulationGroup":      in.OtherPopulationGroup,
		"requestDescription":        in.RequestDescription,
		"responseChannel":           in.ResponseChannel,
		"requestType":               in.RequestType,
		"requestAgainstStudent":     in.RequestAgainstStudent,
		"responsibleStudentName":    in.ResponsibleStudentName,
		"responsibleStudentProgram": in.ResponsibleStudentProgram,
		"studentCaseDescription":    in.StudentCaseDescription,
		"acceptsDataProcessing":     in.AcceptsDataProcessing,
		"submittedAt":               in.SubmittedAt,
		"attachments":               attachmentsPayload,
	}

	bytes, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func internalMailboxRecipients() []string {
	raw := strings.TrimSpace(os.Getenv("PQRS_INTERNAL_EMAIL"))
	if raw == "" {
		raw = strings.TrimSpace(os.Getenv("SMTP_FROM_ADDRESS"))
	}
	if raw == "" {
		raw = strings.TrimSpace(os.Getenv("SMTP_USERNAME"))
	}
	if raw == "" {
		return nil
	}

	parts := strings.Split(raw, ",")
	out := make([]string, 0, len(parts))
	for _, part := range parts {
		value := strings.TrimSpace(part)
		if value == "" {
			continue
		}
		out = append(out, value)
	}
	return out
}

func buildCitizenFullName(firstName string, lastName string) string {
	fullName := strings.TrimSpace(strings.Join([]string{
		strings.TrimSpace(firstName),
		strings.TrimSpace(lastName),
	}, " "))
	if fullName == "" {
		return "citizen"
	}
	return fullName
}
