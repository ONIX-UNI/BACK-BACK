package service

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/DuvanRozoParra/sicou/internal/auth/dto"
)

var errMissingTokenSecret = errors.New("missing token secret")
var errInvalidToken = errors.New("invalid token")
var errExpiredToken = errors.New("expired token")

type tokenClaims struct {
	Subject     string   `json:"sub"`
	Email       string   `json:"email"`
	TokenUse    string   `json:"token_use,omitempty"`
	DisplayName string   `json:"display_name"`
	Roles       []string `json:"roles"`
	IssuedAt    int64    `json:"iat"`
	ExpiresAt   int64    `json:"exp"`
}

type passwordResetTokenClaims struct {
	Subject   string `json:"sub"`
	Email     string `json:"email"`
	TokenUse  string `json:"token_use"`
	IssuedAt  int64  `json:"iat"`
	ExpiresAt int64  `json:"exp"`
}

func buildAccessToken(
	user dto.AppUser,
	roles []string,
	secret string,
	tokenTTL time.Duration,
) (string, time.Time, error) {
	normalizedSecret := strings.TrimSpace(secret)
	if normalizedSecret == "" {
		return "", time.Time{}, errMissingTokenSecret
	}
	if tokenTTL <= 0 {
		tokenTTL = 2 * time.Hour
	}

	now := time.Now().UTC()
	expiresAt := now.Add(tokenTTL)

	header, err := json.Marshal(map[string]string{
		"alg": "HS256",
		"typ": "JWT",
	})
	if err != nil {
		return "", time.Time{}, err
	}

	claims, err := json.Marshal(tokenClaims{
		Subject:     user.ID.String(),
		Email:       user.Email,
		TokenUse:    "session",
		DisplayName: user.DisplayName,
		Roles:       roles,
		IssuedAt:    now.Unix(),
		ExpiresAt:   expiresAt.Unix(),
	})
	if err != nil {
		return "", time.Time{}, err
	}

	encodedHeader := base64.RawURLEncoding.EncodeToString(header)
	encodedClaims := base64.RawURLEncoding.EncodeToString(claims)
	signingInput := encodedHeader + "." + encodedClaims

	signatureMac := hmac.New(sha256.New, []byte(normalizedSecret))
	signatureMac.Write([]byte(signingInput))
	signature := base64.RawURLEncoding.EncodeToString(signatureMac.Sum(nil))

	return signingInput + "." + signature, expiresAt, nil
}

func parseAccessToken(rawToken string, secret string) (tokenClaims, string, time.Time, error) {
	return parseAccessTokenWithExpiry(rawToken, secret, true)
}

func parseAccessTokenAllowExpired(rawToken string, secret string) (tokenClaims, string, time.Time, error) {
	return parseAccessTokenWithExpiry(rawToken, secret, false)
}

func parsePasswordResetToken(rawToken string, secret string) (tokenClaims, string, time.Time, error) {
	claims, signature, expiresAt, err := parseAccessToken(rawToken, secret)
	if err != nil {
		return tokenClaims{}, "", time.Time{}, err
	}

	if !strings.EqualFold(strings.TrimSpace(claims.TokenUse), "password_reset") {
		return tokenClaims{}, "", time.Time{}, errInvalidToken
	}
	if strings.TrimSpace(claims.Email) == "" {
		return tokenClaims{}, "", time.Time{}, errInvalidToken
	}

	return claims, signature, expiresAt, nil
}

func buildPasswordResetToken(
	user dto.AppUser,
	secret string,
	tokenTTL time.Duration,
) (string, time.Time, error) {
	normalizedSecret := strings.TrimSpace(secret)
	if normalizedSecret == "" {
		return "", time.Time{}, errMissingTokenSecret
	}
	if tokenTTL <= 0 {
		tokenTTL = 30 * time.Minute
	}

	now := time.Now().UTC()
	expiresAt := now.Add(tokenTTL)

	header, err := json.Marshal(map[string]string{
		"alg": "HS256",
		"typ": "JWT",
	})
	if err != nil {
		return "", time.Time{}, err
	}

	claims, err := json.Marshal(passwordResetTokenClaims{
		Subject:   user.ID.String(),
		Email:     strings.TrimSpace(user.Email),
		TokenUse:  "password_reset",
		IssuedAt:  now.Unix(),
		ExpiresAt: expiresAt.Unix(),
	})
	if err != nil {
		return "", time.Time{}, err
	}

	encodedHeader := base64.RawURLEncoding.EncodeToString(header)
	encodedClaims := base64.RawURLEncoding.EncodeToString(claims)
	signingInput := encodedHeader + "." + encodedClaims

	signatureMac := hmac.New(sha256.New, []byte(normalizedSecret))
	signatureMac.Write([]byte(signingInput))
	signature := base64.RawURLEncoding.EncodeToString(signatureMac.Sum(nil))

	return signingInput + "." + signature, expiresAt, nil
}

func parseAccessTokenWithExpiry(
	rawToken string,
	secret string,
	checkExpiry bool,
) (tokenClaims, string, time.Time, error) {
	normalizedSecret := strings.TrimSpace(secret)
	if normalizedSecret == "" {
		return tokenClaims{}, "", time.Time{}, errMissingTokenSecret
	}

	normalizedToken := strings.TrimSpace(rawToken)
	if normalizedToken == "" {
		return tokenClaims{}, "", time.Time{}, errInvalidToken
	}

	parts := strings.Split(normalizedToken, ".")
	if len(parts) != 3 {
		return tokenClaims{}, "", time.Time{}, errInvalidToken
	}

	signingInput := parts[0] + "." + parts[1]

	signatureMac := hmac.New(sha256.New, []byte(normalizedSecret))
	signatureMac.Write([]byte(signingInput))
	expectedSignature := base64.RawURLEncoding.EncodeToString(signatureMac.Sum(nil))
	if !hmac.Equal([]byte(expectedSignature), []byte(parts[2])) {
		return tokenClaims{}, "", time.Time{}, errInvalidToken
	}

	decodedClaims, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return tokenClaims{}, "", time.Time{}, errInvalidToken
	}

	var claims tokenClaims
	if err := json.Unmarshal(decodedClaims, &claims); err != nil {
		return tokenClaims{}, "", time.Time{}, errInvalidToken
	}

	if strings.TrimSpace(claims.Subject) == "" || claims.ExpiresAt == 0 {
		return tokenClaims{}, "", time.Time{}, errInvalidToken
	}

	expiresAt := time.Unix(claims.ExpiresAt, 0).UTC()
	if checkExpiry && !expiresAt.After(time.Now().UTC()) {
		return tokenClaims{}, "", time.Time{}, errExpiredToken
	}

	return claims, parts[2], expiresAt, nil
}
