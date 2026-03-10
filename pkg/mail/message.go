package mail

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"mime"
	"mime/quotedprintable"
	stdmail "net/mail"
	"sort"
	"strings"
	"time"
)

type Message struct {
	To          []string
	Cc          []string
	Bcc         []string
	Subject     string
	TextBody    string
	HTMLBody    string
	Attachments []Attachment
	Headers     map[string]string
}

type Attachment struct {
	Filename    string
	ContentType string
	Data        []byte
}

func collectRecipients(msg Message) ([]string, error) {
	rawRecipients := append(append(append([]string{}, msg.To...), msg.Cc...), msg.Bcc...)
	if len(rawRecipients) == 0 {
		return nil, errors.New("at least one recipient is required")
	}

	unique := make(map[string]string, len(rawRecipients))
	for _, raw := range rawRecipients {
		addr, err := parseAddress(raw)
		if err != nil {
			return nil, fmt.Errorf("invalid recipient %q: %w", raw, err)
		}
		unique[strings.ToLower(addr)] = addr
	}

	recipients := make([]string, 0, len(unique))
	for _, addr := range unique {
		recipients = append(recipients, addr)
	}
	sort.Strings(recipients)
	return recipients, nil
}

func buildMessage(cfg SMTPConfig, msg Message) ([]byte, error) {
	if strings.TrimSpace(msg.TextBody) == "" && strings.TrimSpace(msg.HTMLBody) == "" {
		return nil, errors.New("email body is required")
	}

	fromHeader, err := formatAddress(cfg.FromName, cfg.FromAddress)
	if err != nil {
		return nil, fmt.Errorf("invalid from header: %w", err)
	}

	toHeader, err := formatAddressList(msg.To)
	if err != nil {
		return nil, fmt.Errorf("invalid to header: %w", err)
	}
	ccHeader, err := formatAddressList(msg.Cc)
	if err != nil {
		return nil, fmt.Errorf("invalid cc header: %w", err)
	}

	contentType, body, err := buildBody(msg)
	if err != nil {
		return nil, err
	}

	headers := map[string]string{
		"Date":         time.Now().UTC().Format(time.RFC1123Z),
		"From":         fromHeader,
		"Message-ID":   messageID(cfg.FromAddress),
		"MIME-Version": "1.0",
		"Subject":      mime.QEncoding.Encode("utf-8", strings.TrimSpace(msg.Subject)),
	}

	if toHeader != "" {
		headers["To"] = toHeader
	}
	if ccHeader != "" {
		headers["Cc"] = ccHeader
	}

	isMultipart := strings.HasPrefix(contentType, "multipart/")
	headers["Content-Type"] = contentType
	if !isMultipart {
		headers["Content-Transfer-Encoding"] = "quoted-printable"
	}

	for key, value := range msg.Headers {
		key = strings.TrimSpace(key)
		value = strings.TrimSpace(value)
		if key == "" || value == "" {
			continue
		}
		if _, exists := headers[key]; exists {
			continue
		}
		headers[key] = value
	}

	var keys []string
	for key := range headers {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	var buff bytes.Buffer
	for _, key := range keys {
		buff.WriteString(key)
		buff.WriteString(": ")
		buff.WriteString(headers[key])
		buff.WriteString("\r\n")
	}
	buff.WriteString("\r\n")
	buff.Write(body)
	if !bytes.HasSuffix(body, []byte("\r\n")) {
		buff.WriteString("\r\n")
	}

	return buff.Bytes(), nil
}

func buildBody(msg Message) (string, []byte, error) {
	textBody := strings.TrimSpace(msg.TextBody)
	htmlBody := strings.TrimSpace(msg.HTMLBody)

	if textBody != "" && htmlBody != "" {
		if len(msg.Attachments) > 0 {
			return buildMultipartMixedBody(textBody, htmlBody, msg.Attachments)
		}

		boundary := fmt.Sprintf("b-%d", time.Now().UnixNano())
		var buff bytes.Buffer

		if err := writeMultipartPart(&buff, boundary, "text/plain; charset=utf-8", textBody); err != nil {
			return "", nil, err
		}
		if err := writeMultipartPart(&buff, boundary, "text/html; charset=utf-8", htmlBody); err != nil {
			return "", nil, err
		}
		buff.WriteString("--")
		buff.WriteString(boundary)
		buff.WriteString("--\r\n")

		return fmt.Sprintf("multipart/alternative; boundary=%q", boundary), buff.Bytes(), nil
	}

	if htmlBody != "" {
		if len(msg.Attachments) > 0 {
			return buildMultipartMixedBody("", htmlBody, msg.Attachments)
		}

		encoded, err := quotedPrintable(htmlBody)
		if err != nil {
			return "", nil, err
		}
		return "text/html; charset=utf-8", encoded, nil
	}

	if len(msg.Attachments) > 0 {
		return buildMultipartMixedBody(textBody, "", msg.Attachments)
	}

	encoded, err := quotedPrintable(textBody)
	if err != nil {
		return "", nil, err
	}
	return "text/plain; charset=utf-8", encoded, nil
}

func buildMultipartMixedBody(textBody string, htmlBody string, attachments []Attachment) (string, []byte, error) {
	mixedBoundary := fmt.Sprintf("mixed-%d", time.Now().UnixNano())
	var buff bytes.Buffer

	switch {
	case textBody != "" && htmlBody != "":
		altBoundary := fmt.Sprintf("alt-%d", time.Now().UnixNano())
		buff.WriteString("--")
		buff.WriteString(mixedBoundary)
		buff.WriteString("\r\n")
		buff.WriteString("Content-Type: multipart/alternative; boundary=\"")
		buff.WriteString(altBoundary)
		buff.WriteString("\"\r\n\r\n")

		if err := writeMultipartPart(&buff, altBoundary, "text/plain; charset=utf-8", textBody); err != nil {
			return "", nil, err
		}
		if err := writeMultipartPart(&buff, altBoundary, "text/html; charset=utf-8", htmlBody); err != nil {
			return "", nil, err
		}

		buff.WriteString("--")
		buff.WriteString(altBoundary)
		buff.WriteString("--\r\n")
	case htmlBody != "":
		if err := writeMultipartPart(&buff, mixedBoundary, "text/html; charset=utf-8", htmlBody); err != nil {
			return "", nil, err
		}
	default:
		if err := writeMultipartPart(&buff, mixedBoundary, "text/plain; charset=utf-8", textBody); err != nil {
			return "", nil, err
		}
	}

	for _, attachment := range attachments {
		filename := strings.TrimSpace(attachment.Filename)
		if filename == "" {
			filename = "attachment.bin"
		}

		contentType := strings.TrimSpace(attachment.ContentType)
		if contentType == "" {
			contentType = "application/octet-stream"
		}

		buff.WriteString("--")
		buff.WriteString(mixedBoundary)
		buff.WriteString("\r\n")
		buff.WriteString("Content-Type: ")
		buff.WriteString(contentType)
		buff.WriteString("\r\n")
		buff.WriteString("Content-Transfer-Encoding: base64\r\n")
		buff.WriteString("Content-Disposition: attachment; filename=\"")
		buff.WriteString(escapeHeaderParam(filename))
		buff.WriteString("\"\r\n\r\n")
		buff.Write(base64MIMEEncode(attachment.Data))
	}

	buff.WriteString("--")
	buff.WriteString(mixedBoundary)
	buff.WriteString("--\r\n")

	return fmt.Sprintf("multipart/mixed; boundary=%q", mixedBoundary), buff.Bytes(), nil
}

func writeMultipartPart(buff *bytes.Buffer, boundary string, contentType string, content string) error {
	encoded, err := quotedPrintable(content)
	if err != nil {
		return err
	}
	buff.WriteString("--")
	buff.WriteString(boundary)
	buff.WriteString("\r\n")
	buff.WriteString("Content-Type: ")
	buff.WriteString(contentType)
	buff.WriteString("\r\n")
	buff.WriteString("Content-Transfer-Encoding: quoted-printable\r\n\r\n")
	buff.Write(encoded)
	if !bytes.HasSuffix(encoded, []byte("\r\n")) {
		buff.WriteString("\r\n")
	}

	return nil
}

func quotedPrintable(value string) ([]byte, error) {
	var buff bytes.Buffer
	writer := quotedprintable.NewWriter(&buff)
	if _, err := writer.Write([]byte(normalizeLineBreaks(value))); err != nil {
		return nil, err
	}
	if err := writer.Close(); err != nil {
		return nil, err
	}
	return buff.Bytes(), nil
}

func base64MIMEEncode(data []byte) []byte {
	encoded := base64.StdEncoding.EncodeToString(data)
	if encoded == "" {
		return []byte("\r\n")
	}

	var buff bytes.Buffer
	for len(encoded) > 76 {
		buff.WriteString(encoded[:76])
		buff.WriteString("\r\n")
		encoded = encoded[76:]
	}
	if len(encoded) > 0 {
		buff.WriteString(encoded)
	}
	buff.WriteString("\r\n")
	return buff.Bytes()
}

func escapeHeaderParam(value string) string {
	value = strings.ReplaceAll(value, "\\", "\\\\")
	value = strings.ReplaceAll(value, "\"", "\\\"")
	return value
}

func normalizeLineBreaks(value string) string {
	value = strings.ReplaceAll(value, "\r\n", "\n")
	value = strings.ReplaceAll(value, "\r", "\n")
	return strings.ReplaceAll(value, "\n", "\r\n")
}

func formatAddress(name string, address string) (string, error) {
	parsed, err := stdmail.ParseAddress(strings.TrimSpace(address))
	if err != nil {
		return "", err
	}

	if strings.TrimSpace(name) == "" {
		return parsed.Address, nil
	}

	return (&stdmail.Address{
		Name:    strings.TrimSpace(name),
		Address: parsed.Address,
	}).String(), nil
}

func formatAddressList(values []string) (string, error) {
	if len(values) == 0 {
		return "", nil
	}

	addresses := make([]string, 0, len(values))
	for _, value := range values {
		address, err := parseAddress(value)
		if err != nil {
			return "", err
		}
		addresses = append(addresses, address)
	}

	return strings.Join(addresses, ", "), nil
}

func parseAddress(value string) (string, error) {
	parsed, err := stdmail.ParseAddress(strings.TrimSpace(value))
	if err != nil {
		return "", err
	}
	return parsed.Address, nil
}

func messageID(fromAddress string) string {
	parts := strings.Split(fromAddress, "@")
	domain := "localhost"
	if len(parts) == 2 && strings.TrimSpace(parts[1]) != "" {
		domain = strings.TrimSpace(parts[1])
	}
	return fmt.Sprintf("<%d.%d@%s>", time.Now().UnixNano(), time.Now().Unix(), domain)
}
