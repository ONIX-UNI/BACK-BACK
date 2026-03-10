package dto

import "strings"

func NormalizeRoleCode(role string) string {
	return strings.ToUpper(strings.TrimSpace(role))
}

func NormalizeAuditAction(action string) string {
	return strings.ToUpper(strings.TrimSpace(action))
}

func NormalizeRoleCodes(roles []string) []string {
	seen := make(map[string]struct{}, len(roles))
	normalized := make([]string, 0, len(roles))

	for _, role := range roles {
		value := NormalizeRoleCode(role)
		if value == "" {
			continue
		}
		if _, exists := seen[value]; exists {
			continue
		}
		seen[value] = struct{}{}
		normalized = append(normalized, value)
	}

	return normalized
}
