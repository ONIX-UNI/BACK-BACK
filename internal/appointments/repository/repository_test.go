package repository

import (
	"strings"
	"testing"
	"time"
)

func TestFormatPreturnoNumber(t *testing.T) {
	createdAt := time.Date(2026, time.February, 16, 13, 59, 9, 0, time.UTC)

	got := formatPreturnoNumber(createdAt, 31)
	want := "PT-2026-0031"
	if got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}

func TestBuildCitizenNotificationBodyIncludesPreturno(t *testing.T) {
	body := buildCitizenNotificationBody("PT-2026-0031", "Daniel Montillo")

	expectedFragments := []string{
		"Hola Daniel Montillo,",
		"Gracias por registrar tu solicitud de asesoria juridica.",
		"Tu numero de preturno es PT-2026-0031.",
	}
	for _, fragment := range expectedFragments {
		if !strings.Contains(body, fragment) {
			t.Fatalf("expected body to include %q, got:\n%s", fragment, body)
		}
	}
}
