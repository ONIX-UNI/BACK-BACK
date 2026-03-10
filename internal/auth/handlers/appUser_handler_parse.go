package handlers

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

func parsePositiveInt(value string, defaultValue int) int {
	raw := strings.TrimSpace(value)
	if raw == "" {
		return defaultValue
	}

	parsed, err := strconv.Atoi(raw)
	if err != nil || parsed <= 0 {
		return defaultValue
	}

	return parsed
}

func parseStrictPositiveInt(value string, defaultValue int) (int, error) {
	raw := strings.TrimSpace(value)
	if raw == "" {
		return defaultValue, nil
	}

	parsed, err := strconv.Atoi(raw)
	if err != nil || parsed <= 0 {
		return 0, errors.New("invalid positive integer")
	}

	return parsed, nil
}

func parseNonNegativeInt(value string, defaultValue int) int {
	raw := strings.TrimSpace(value)
	if raw == "" {
		return defaultValue
	}

	parsed, err := strconv.Atoi(raw)
	if err != nil || parsed < 0 {
		return defaultValue
	}

	return parsed
}

func parseStrictNonNegativeInt(value string) (int, error) {
	raw := strings.TrimSpace(value)
	if raw == "" {
		return 0, errors.New("value is required")
	}

	parsed, err := strconv.Atoi(raw)
	if err != nil || parsed < 0 {
		return 0, errors.New("invalid non-negative integer")
	}

	return parsed, nil
}

func parseLimit(value string, defaultValue, max int) int {
	limit := parsePositiveInt(value, defaultValue)
	if limit > max {
		limit = max
	}
	return limit
}

func parseStrictLimit(value string, defaultValue, max int) (int, error) {
	limit, err := parseStrictPositiveInt(value, defaultValue)
	if err != nil {
		return 0, err
	}
	if limit > max {
		limit = max
	}
	return limit, nil
}

func parseOptionalBool(value string, defaultValue bool) (bool, error) {
	raw := strings.TrimSpace(value)
	if raw == "" {
		return defaultValue, nil
	}

	parsed, err := strconv.ParseBool(raw)
	if err != nil {
		return defaultValue, err
	}

	return parsed, nil
}

func parseISODatetime(value string) (time.Time, error) {
	layouts := []string{
		time.RFC3339Nano,
		time.RFC3339,
	}

	for _, layout := range layouts {
		if parsed, err := time.Parse(layout, value); err == nil {
			return parsed, nil
		}
	}

	return time.Time{}, errors.New("invalid datetime")
}
