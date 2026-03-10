package mail

import (
	"strings"
	"testing"
	"time"
)

func validGoogleConfig() SMTPConfig {
	return SMTPConfig{
		Host:           "smtp.gmail.com",
		Port:           587,
		Username:       "noreply@institucion.edu.co",
		Password:       "app-password",
		AuthMethod:     AuthMethodPlain,
		FromAddress:    "noreply@institucion.edu.co",
		UseStartTLS:    true,
		RequireTLS:     true,
		Timeout:        10 * time.Second,
		FromName:       "SICOU",
		UseImplicitTLS: false,
	}
}

func TestSMTPConfigValidateGoogleOK(t *testing.T) {
	cfg := validGoogleConfig()
	if err := cfg.Validate(); err != nil {
		t.Fatalf("expected config to be valid, got error: %v", err)
	}
}

func TestSMTPConfigValidateGoogleWithoutTLSFails(t *testing.T) {
	cfg := validGoogleConfig()
	cfg.RequireTLS = false
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected tls validation error for google smtp")
	}
}

func TestSMTPConfigValidateXOAuth2(t *testing.T) {
	cfg := validGoogleConfig()
	cfg.AuthMethod = AuthMethodXOAuth2
	cfg.Password = ""
	cfg.AccessToken = "token"
	if err := cfg.Validate(); err != nil {
		t.Fatalf("expected xoauth2 config to be valid, got error: %v", err)
	}
}

func TestCollectRecipientsDeduplicates(t *testing.T) {
	msg := Message{
		To:  []string{"foo@example.com"},
		Cc:  []string{"Foo@example.com"},
		Bcc: []string{"bar@example.com"},
	}

	recipients, err := collectRecipients(msg)
	if err != nil {
		t.Fatalf("collectRecipients returned error: %v", err)
	}

	if len(recipients) != 2 {
		t.Fatalf("expected 2 unique recipients, got %d", len(recipients))
	}
}

func TestBuildMessageMultipart(t *testing.T) {
	cfg := validGoogleConfig()
	msg := Message{
		To:       []string{"destino@example.com"},
		Subject:  "Prueba",
		TextBody: "texto",
		HTMLBody: "<strong>html</strong>",
	}

	payload, err := buildMessage(cfg, msg)
	if err != nil {
		t.Fatalf("buildMessage returned error: %v", err)
	}

	raw := string(payload)
	if !strings.Contains(raw, "Content-Type: multipart/alternative") {
		t.Fatal("expected multipart content type")
	}
	if !strings.Contains(raw, "Content-Transfer-Encoding: quoted-printable") {
		t.Fatal("expected quoted printable encoding")
	}
	if !strings.Contains(raw, "To: destino@example.com") {
		t.Fatal("expected To header")
	}
}

func TestBuildMessageWithAttachmentUsesMultipartMixed(t *testing.T) {
	cfg := validGoogleConfig()
	msg := Message{
		To:       []string{"destino@example.com"},
		Subject:  "Prueba con adjunto",
		TextBody: "texto",
		Attachments: []Attachment{
			{
				Filename:    "evidencia.txt",
				ContentType: "text/plain",
				Data:        []byte("contenido"),
			},
		},
	}

	payload, err := buildMessage(cfg, msg)
	if err != nil {
		t.Fatalf("buildMessage returned error: %v", err)
	}

	raw := string(payload)
	if !strings.Contains(raw, "Content-Type: multipart/mixed") {
		t.Fatal("expected multipart mixed content type")
	}
	if !strings.Contains(raw, "Content-Disposition: attachment; filename=\"evidencia.txt\"") {
		t.Fatal("expected attachment disposition header")
	}
	if !strings.Contains(raw, "Content-Transfer-Encoding: base64") {
		t.Fatal("expected base64 encoding for attachment")
	}
}
