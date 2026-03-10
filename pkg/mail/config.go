package mail

import (
	"errors"
	"fmt"
	stdmail "net/mail"
	"strings"
	"time"
)

const (
	defaultSMTPHost    = "smtp.gmail.com"
	defaultSMTPPort    = 587
	defaultSMTPTimeout = 10 * time.Second
)

type AuthMethod string

const (
	AuthMethodNone    AuthMethod = "none"
	AuthMethodPlain   AuthMethod = "plain"
	AuthMethodXOAuth2 AuthMethod = "xoauth2"
)

type SMTPConfig struct {
	Host               string
	Port               int
	Username           string
	Password           string
	AccessToken        string
	AuthMethod         AuthMethod
	FromName           string
	FromAddress        string
	UseStartTLS        bool
	UseImplicitTLS     bool
	RequireTLS         bool
	InsecureSkipVerify bool
	Timeout            time.Duration
}

func (c *SMTPConfig) Validate() error {
	host := strings.TrimSpace(c.Host)
	if host == "" {
		return errors.New("smtp host is required")
	}

	if c.Port <= 0 || c.Port > 65535 {
		return errors.New("smtp port must be between 1 and 65535")
	}

	if c.UseStartTLS && c.UseImplicitTLS {
		return errors.New("use_starttls and use_implicit_tls are mutually exclusive")
	}

	if c.RequireTLS && !c.UseStartTLS && !c.UseImplicitTLS {
		return errors.New("require_tls requires starttls or implicit tls")
	}

	if c.Timeout <= 0 {
		return errors.New("smtp timeout must be greater than zero")
	}

	if _, err := stdmail.ParseAddress(strings.TrimSpace(c.FromAddress)); err != nil {
		return fmt.Errorf("invalid smtp from address: %w", err)
	}

	authMethod := normalizeAuthMethod(c.AuthMethod)
	switch authMethod {
	case AuthMethodNone:
		// SMTP relay servers may allow no authentication.
	case AuthMethodPlain:
		if strings.TrimSpace(c.Username) == "" || strings.TrimSpace(c.Password) == "" {
			return errors.New("plain auth requires smtp username and smtp password")
		}
	case AuthMethodXOAuth2:
		if strings.TrimSpace(c.Username) == "" || strings.TrimSpace(c.AccessToken) == "" {
			return errors.New("xoauth2 auth requires smtp username and smtp access token")
		}
	default:
		return fmt.Errorf("unsupported smtp auth method %q", c.AuthMethod)
	}
	c.AuthMethod = authMethod

	if isGoogleSMTPHost(host) {
		if c.Port != 587 && c.Port != 465 {
			return errors.New("google smtp requires port 587 (starttls) or 465 (implicit tls)")
		}

		if !c.RequireTLS {
			return errors.New("google smtp requires tls")
		}

		if c.Port == 587 && !c.UseStartTLS {
			return errors.New("google smtp on port 587 requires starttls")
		}

		if c.Port == 465 && !c.UseImplicitTLS {
			return errors.New("google smtp on port 465 requires implicit tls")
		}
	}

	return nil
}

func withDefaults(cfg SMTPConfig) SMTPConfig {
	cfg.Host = strings.TrimSpace(cfg.Host)
	if cfg.Host == "" {
		cfg.Host = defaultSMTPHost
	}

	if cfg.Port == 0 {
		cfg.Port = defaultSMTPPort
	}

	if cfg.Timeout <= 0 {
		cfg.Timeout = defaultSMTPTimeout
	}

	if cfg.AuthMethod == "" {
		cfg.AuthMethod = AuthMethodPlain
	}

	if cfg.RequireTLS && cfg.Port == 465 && !cfg.UseImplicitTLS && !cfg.UseStartTLS {
		cfg.UseImplicitTLS = true
	}
	if cfg.RequireTLS && cfg.Port == 587 && !cfg.UseImplicitTLS && !cfg.UseStartTLS {
		cfg.UseStartTLS = true
	}

	cfg.FromAddress = strings.TrimSpace(cfg.FromAddress)
	if cfg.FromAddress == "" {
		cfg.FromAddress = strings.TrimSpace(cfg.Username)
	}

	return cfg
}

func normalizeAuthMethod(value AuthMethod) AuthMethod {
	normalized := AuthMethod(strings.ToLower(strings.TrimSpace(string(value))))
	if normalized == "" {
		return AuthMethodPlain
	}

	switch normalized {
	case AuthMethodNone:
		return AuthMethodNone
	case AuthMethodPlain:
		return AuthMethodPlain
	case AuthMethodXOAuth2:
		return AuthMethodXOAuth2
	default:
		return normalized
	}
}

func isGoogleSMTPHost(host string) bool {
	host = strings.ToLower(strings.TrimSpace(host))
	return strings.HasSuffix(host, "gmail.com") || strings.HasSuffix(host, "googlemail.com")
}
