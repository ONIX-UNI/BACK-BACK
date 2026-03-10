package service

import "testing"

func TestNormalizeAndValidate_AllowsAnonymousQueryTypeAsYes(t *testing.T) {
	in := baseInput()
	in.QueryType = "Si"

	_, err := normalizeAndValidate(in)
	if err != nil {
		t.Fatalf("expected validation success for anonymous queryType, got error: %+v", err.Fields)
	}
}

func TestNormalizeAndValidate_QueryTypeNoRequiresIdentityFields(t *testing.T) {
	in := baseInput()
	in.QueryType = "No"

	_, err := normalizeAndValidate(in)
	if err == nil {
		t.Fatalf("expected validation error when queryType is identified and identity fields are empty")
	}

	requiredFields := []string{
		"personType",
		"documentType",
		"documentOrTaxId",
		"firstName",
		"firstLastName",
		"gender",
		"address",
		"neighborhoodArea",
		"email",
		"phone",
	}
	for _, field := range requiredFields {
		if _, ok := err.Fields[field]; !ok {
			t.Fatalf("expected validation error for field %q, got: %+v", field, err.Fields)
		}
	}
}

func TestNormalizeAndValidate_InvalidQueryTypeDoesNotCascadeIdentityErrors(t *testing.T) {
	in := baseInput()
	in.QueryType = "Tal vez"

	_, err := normalizeAndValidate(in)
	if err == nil {
		t.Fatalf("expected validation error for invalid queryType")
	}

	if got := err.Fields["queryType"]; got == "" {
		t.Fatalf("expected queryType validation error, got: %+v", err.Fields)
	}

	if len(err.Fields) != 1 {
		t.Fatalf("expected only queryType validation error, got: %+v", err.Fields)
	}
}

func baseInput() CreatePQRSInput {
	return CreatePQRSInput{
		AllowsElectronicResponse: "No",
		RequestAgainstStudent:    "No",
		AcceptsDataProcessing:    "Si",
		RequestDescription:       "Descripcion de prueba",
		ResponseChannel:          "No aplica",
		RequestType:              "Peticion",
		SubmittedAt:              "2026-02-16T16:16:18.218Z",
	}
}
