package repository

import (
	"strings"
	"testing"
	"time"
)

func TestFormatRadicado(t *testing.T) {
	createdAt := time.Date(2026, time.February, 16, 11, 16, 18, 0, time.UTC)

	got := formatRadicado(createdAt, 45)
	want := "PQRS-2026-0045"
	if got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}

func TestBuildInternalEmailBody_IncludesPayloadFields(t *testing.T) {
	bodyPayload := `{
		"queryType":"No",
		"documentType":"Cedula de ciudadania",
		"documentOrTaxId":"1006825436",
		"requestType":"Queja",
		"requestDescription":"No atendieron mi caso",
		"responseChannel":"No aplica",
		"firstName":"Ana",
		"middleName":"Maria",
		"firstLastName":"Perez",
		"secondLastName":"Gomez",
		"requestAgainstStudent":"No",
		"email":"ana@example.com",
		"phone":"3001234567",
		"submittedAt":"2026-02-16T16:16:18.218Z",
		"attachments":[
			{"name":"evidencia.pdf","type":"application/pdf","size":12345,"url":"https://files/evidencia.pdf"}
		]
	}`

	got := buildInternalEmailBody(
		"PQRS-2026-0006",
		"",
		"citizen",
		nil,
		bodyPayload,
	)

	expectedFragments := []string{
		"Se recibio una nueva PQRS en SICOU.",
		"Radicado: PQRS-2026-0006",
		"Tipo de solicitud: Queja",
		"Asunto: No atendieron mi caso",
		"Canal de respuesta preferido: No aplica",
		"- Nombre completo: Ana Maria Perez Gomez",
		"- Documento: Cedula de ciudadania 1006825436",
		"- Correo electronico: ana@example.com",
		"Descripcion de la solicitud:\nNo atendieron mi caso",
		"Informacion sobre posible estudiante involucrado:",
		"- Aplica: No",
		"Adjuntos:\n- evidencia.pdf",
	}

	for _, fragment := range expectedFragments {
		if !strings.Contains(got, fragment) {
			t.Fatalf("expected email body to include %q, got body:\n%s", fragment, got)
		}
	}
}

func TestBuildInternalEmailBody_InvalidPayloadKeepsBaseSummary(t *testing.T) {
	got := buildInternalEmailBody(
		"PQRS-2026-0007",
		"Asunto base",
		"",
		nil,
		"{invalid-json",
	)

	if !strings.Contains(got, "Radicado: PQRS-2026-0007") {
		t.Fatalf("expected radicado in body, got:\n%s", got)
	}
	if !strings.Contains(got, "Asunto: Asunto base") {
		t.Fatalf("expected asunto in body, got:\n%s", got)
	}
	if !strings.Contains(got, "Se recibio una nueva PQRS en SICOU.") {
		t.Fatalf("expected human-readable header in body, got:\n%s", got)
	}
	if strings.Contains(got, "queryType:") {
		t.Fatalf("did not expect technical keys in body, got:\n%s", got)
	}
}
