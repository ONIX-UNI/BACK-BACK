package service

import (
	"testing"
	"time"
)

func TestNormalizeAndValidate_MinimumValidInput(t *testing.T) {
	in := CreateAsesoriaInput{
		AcceptsDataProcessing: "Si",
		ConsultationDate:      "2026-02-16",
		FullName:              "Daniel Montillo Cardenas",
		DocumentType:          "Cedula de ciudadania",
		DocumentNumber:        "1006825436",
		CaseDescription:       "Descripcion de caso de prueba",
		SubmittedAt:           "2026-02-16T16:49:46.076Z",
	}

	normalized, err := normalizeAndValidate(in)
	if err != nil {
		t.Fatalf("expected validation success, got error: %+v", err.Fields)
	}

	if !normalized.AcceptsDataProcessingBool {
		t.Fatalf("expected acceptsDataProcessing to be true")
	}
}

func TestNormalizeAndValidate_RequiredFields(t *testing.T) {
	_, err := normalizeAndValidate(CreateAsesoriaInput{})
	if err == nil {
		t.Fatalf("expected validation error for empty input")
	}

	required := []string{
		"acceptsDataProcessing",
		"consultationDate",
		"fullName",
		"documentType",
		"documentNumber",
		"caseDescription",
		"submittedAt",
	}

	for _, field := range required {
		if _, ok := err.Fields[field]; !ok {
			t.Fatalf("expected validation error for %q, got: %+v", field, err.Fields)
		}
	}
}

func TestNormalizeAndValidate_AuthorizesNotificationRequiresEmail(t *testing.T) {
	in := CreateAsesoriaInput{
		AcceptsDataProcessing:  "Si",
		ConsultationDate:       "2026-02-16",
		FullName:               "Daniel Montillo Cardenas",
		DocumentType:           "Cedula de ciudadania",
		DocumentNumber:         "1006825436",
		CaseDescription:        "Descripcion de caso de prueba",
		AuthorizesNotification: "Si",
		SubmittedAt:            "2026-02-16T16:49:46.076Z",
	}

	_, err := normalizeAndValidate(in)
	if err == nil {
		t.Fatalf("expected validation error when authorizesNotification is Yes and email is empty")
	}
	if got := err.Fields["email"]; got == "" {
		t.Fatalf("expected email validation error, got: %+v", err.Fields)
	}
}

func TestSupportLinks_FromFileObjectID(t *testing.T) {
	viewURL, downloadURL := supportLinks("45d3e681-433a-4316-965d-141f1a37662b", "")

	if viewURL != "/api/v1/documents/files/45d3e681-433a-4316-965d-141f1a37662b/view" {
		t.Fatalf("unexpected viewURL: %s", viewURL)
	}
	if downloadURL != "/api/v1/documents/files/45d3e681-433a-4316-965d-141f1a37662b/download" {
		t.Fatalf("unexpected downloadURL: %s", downloadURL)
	}
}

func TestFormatDateOnly_UTCNormalizationForTimestamp(t *testing.T) {
	got := formatDateOnly("2026-02-19T19:00:00-05:00", time.Time{})
	if got != "2026-02-20" {
		t.Fatalf("expected 2026-02-20, got %s", got)
	}
}
