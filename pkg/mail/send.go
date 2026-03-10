package mail

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"net/smtp"
	"time"
)

type Sender interface {
	Send(ctx context.Context, msg Message) error
}

func (c *SMTPClient) Send(ctx context.Context, msg Message) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	recipients, err := collectRecipients(msg)
	if err != nil {
		return err
	}

	payload, err := buildMessage(c.cfg, msg)
	if err != nil {
		return err
	}

	conn, client, err := c.connect(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()
	defer client.Close()

	if err := c.authenticate(client); err != nil {
		return err
	}

	if err := client.Mail(c.cfg.FromAddress); err != nil {
		return fmt.Errorf("smtp MAIL FROM failed: %w", err)
	}

	for _, recipient := range recipients {
		if err := client.Rcpt(recipient); err != nil {
			return fmt.Errorf("smtp RCPT TO failed for %s: %w", recipient, err)
		}
	}

	dataWriter, err := client.Data()
	if err != nil {
		return fmt.Errorf("smtp DATA failed: %w", err)
	}

	if _, err := dataWriter.Write(payload); err != nil {
		_ = dataWriter.Close()
		return fmt.Errorf("smtp write body failed: %w", err)
	}

	if err := dataWriter.Close(); err != nil {
		return fmt.Errorf("smtp DATA close failed: %w", err)
	}

	if err := client.Quit(); err != nil {
		return fmt.Errorf("smtp QUIT failed: %w", err)
	}

	return nil
}

func (c *SMTPClient) connect(ctx context.Context) (net.Conn, *smtp.Client, error) {
	addr := fmt.Sprintf("%s:%d", c.cfg.Host, c.cfg.Port)
	dialer := &net.Dialer{Timeout: c.cfg.Timeout}
	tlsCfg := &tls.Config{
		MinVersion:         tls.VersionTLS12,
		ServerName:         c.cfg.Host,
		InsecureSkipVerify: c.cfg.InsecureSkipVerify, //nolint:gosec // Controlled by config for local/dev only.
	}

	var conn net.Conn
	var err error

	if c.cfg.UseImplicitTLS {
		tlsDialer := &tls.Dialer{
			NetDialer: dialer,
			Config:    tlsCfg,
		}
		conn, err = tlsDialer.DialContext(ctx, "tcp", addr)
	} else {
		conn, err = dialer.DialContext(ctx, "tcp", addr)
	}
	if err != nil {
		return nil, nil, fmt.Errorf("smtp dial failed: %w", err)
	}

	if err := conn.SetDeadline(time.Now().Add(c.cfg.Timeout)); err != nil {
		_ = conn.Close()
		return nil, nil, fmt.Errorf("smtp set deadline failed: %w", err)
	}

	client, err := smtp.NewClient(conn, c.cfg.Host)
	if err != nil {
		_ = conn.Close()
		return nil, nil, fmt.Errorf("smtp client creation failed: %w", err)
	}

	if !c.cfg.UseImplicitTLS && (c.cfg.UseStartTLS || c.cfg.RequireTLS) {
		hasStartTLS, _ := client.Extension("STARTTLS")
		if !hasStartTLS {
			if c.cfg.RequireTLS {
				_ = client.Close()
				_ = conn.Close()
				return nil, nil, errors.New("smtp server does not support STARTTLS")
			}
		} else if err := client.StartTLS(tlsCfg); err != nil {
			_ = client.Close()
			_ = conn.Close()
			return nil, nil, fmt.Errorf("smtp STARTTLS failed: %w", err)
		}
	}

	return conn, client, nil
}

func (c *SMTPClient) authenticate(client *smtp.Client) error {
	method := normalizeAuthMethod(c.cfg.AuthMethod)
	if method == AuthMethodNone {
		return nil
	}

	hasAuth, _ := client.Extension("AUTH")
	if !hasAuth {
		return errors.New("smtp server does not advertise AUTH")
	}

	var auth smtp.Auth
	switch method {
	case AuthMethodPlain:
		auth = smtp.PlainAuth("", c.cfg.Username, c.cfg.Password, c.cfg.Host)
	case AuthMethodXOAuth2:
		auth = xoauth2Auth{
			username:    c.cfg.Username,
			accessToken: c.cfg.AccessToken,
		}
	default:
		return fmt.Errorf("unsupported smtp auth method %q", method)
	}

	if err := client.Auth(auth); err != nil {
		return fmt.Errorf("smtp auth failed: %w", err)
	}

	return nil
}

type xoauth2Auth struct {
	username    string
	accessToken string
}

func (a xoauth2Auth) Start(server *smtp.ServerInfo) (string, []byte, error) {
	if !server.TLS {
		return "", nil, errors.New("xoauth2 authentication requires tls")
	}

	payload := fmt.Sprintf("user=%s\x01auth=Bearer %s\x01\x01", a.username, a.accessToken)
	return "XOAUTH2", []byte(payload), nil
}

func (a xoauth2Auth) Next(fromServer []byte, more bool) ([]byte, error) {
	if more {
		return nil, fmt.Errorf("unexpected xoauth2 challenge: %s", string(fromServer))
	}
	return nil, nil
}
