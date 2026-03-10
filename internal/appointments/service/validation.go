package service

import (
	"fmt"
	stdmail "net/mail"
	"strconv"
	"strings"
	"time"
)

type normalizedInput struct {
	CreateAsesoriaInput
	AcceptsDataProcessingBool  bool
	AuthorizesNotificationBool bool
	HeadOfHouseholdBool        *bool
	ConsultationDateAt         time.Time
	SubmittedAtTime            time.Time
}

func normalizeAndValidate(in CreateAsesoriaInput) (normalizedInput, *ValidationError) {
	fieldsErr := map[string]string{}

	in.AcceptsDataProcessing = strings.TrimSpace(in.AcceptsDataProcessing)
	in.ConsultationDate = strings.TrimSpace(in.ConsultationDate)
	in.FullName = strings.TrimSpace(in.FullName)
	in.DocumentType = strings.TrimSpace(in.DocumentType)
	in.OtherDocumentType = strings.TrimSpace(in.OtherDocumentType)
	in.DocumentNumber = strings.TrimSpace(in.DocumentNumber)
	in.BirthDate = strings.TrimSpace(in.BirthDate)
	in.Age = strings.TrimSpace(in.Age)
	in.MaritalStatus = strings.TrimSpace(in.MaritalStatus)
	in.OtherMaritalStatus = strings.TrimSpace(in.OtherMaritalStatus)
	in.Gender = strings.TrimSpace(in.Gender)
	in.Address = strings.TrimSpace(in.Address)
	in.HousingType = strings.TrimSpace(in.HousingType)
	in.OtherHousingType = strings.TrimSpace(in.OtherHousingType)
	in.SocioEconomicStratum = strings.TrimSpace(in.SocioEconomicStratum)
	in.SisbenCategory = strings.TrimSpace(in.SisbenCategory)
	in.MobilePhone = strings.TrimSpace(in.MobilePhone)
	in.Email = strings.TrimSpace(in.Email)
	in.PopulationType = strings.TrimSpace(in.PopulationType)
	in.OtherPopulationType = strings.TrimSpace(in.OtherPopulationType)
	in.HeadOfHousehold = strings.TrimSpace(in.HeadOfHousehold)
	in.Occupation = strings.TrimSpace(in.Occupation)
	in.EducationLevel = strings.TrimSpace(in.EducationLevel)
	in.OtherEducationLevel = strings.TrimSpace(in.OtherEducationLevel)
	in.CaseDescription = strings.TrimSpace(in.CaseDescription)
	in.AuthorizesNotification = strings.TrimSpace(in.AuthorizesNotification)
	in.SubmittedAt = strings.TrimSpace(in.SubmittedAt)
	in.ActorUserID = strings.TrimSpace(in.ActorUserID)
	in.TimelineSource = strings.TrimSpace(in.TimelineSource)

	acceptsDataProcessing, errAccepts := parseOptionalYesNo(in.AcceptsDataProcessing)
	if errAccepts != nil {
		fieldsErr["acceptsDataProcessing"] = "must be Yes or No"
	}
	if !acceptsDataProcessing {
		fieldsErr["acceptsDataProcessing"] = "must accept data processing"
	}

	if in.FullName == "" {
		fieldsErr["fullName"] = "is required"
	}
	if in.DocumentType == "" {
		fieldsErr["documentType"] = "is required"
	}
	if in.DocumentNumber == "" {
		fieldsErr["documentNumber"] = "is required"
	}
	if in.CaseDescription == "" {
		fieldsErr["caseDescription"] = "is required"
	}
	if in.ConsultationDate == "" {
		fieldsErr["consultationDate"] = "is required"
	}
	if in.SubmittedAt == "" {
		fieldsErr["submittedAt"] = "is required"
	}

	consultationDateAt, errConsultation := parseDateTime(in.ConsultationDate)
	if errConsultation != nil {
		fieldsErr["consultationDate"] = "invalid date format"
	}

	submittedAtTime, errSubmitted := parseDateTime(in.SubmittedAt)
	if errSubmitted != nil {
		fieldsErr["submittedAt"] = "invalid date format"
	}

	authorizesNotification := false
	if in.AuthorizesNotification != "" {
		value, err := parseOptionalYesNo(in.AuthorizesNotification)
		if err != nil {
			fieldsErr["authorizesNotification"] = "must be Yes or No"
		}
		authorizesNotification = value
	}

	if in.Email != "" {
		if _, err := stdmail.ParseAddress(in.Email); err != nil {
			fieldsErr["email"] = "invalid email"
		}
	}
	if authorizesNotification && in.Email == "" {
		fieldsErr["email"] = "is required when authorizesNotification is Yes"
	}

	if in.Age != "" {
		if _, err := parseNonNegativeInt(in.Age); err != nil {
			fieldsErr["age"] = "must be a valid non-negative integer"
		}
	}

	if in.SocioEconomicStratum != "" {
		stratum, err := parseNonNegativeInt(in.SocioEconomicStratum)
		if err != nil || stratum < 1 || stratum > 6 {
			fieldsErr["socioEconomicStratum"] = "must be a number between 1 and 6"
		}
	}

	var headOfHousehold *bool
	if in.HeadOfHousehold != "" {
		value, err := parseOptionalYesNo(in.HeadOfHousehold)
		if err != nil {
			fieldsErr["headOfHousehold"] = "must be Yes or No"
		}
		headOfHousehold = &value
	}

	if len(fieldsErr) > 0 {
		return normalizedInput{}, &ValidationError{Fields: fieldsErr}
	}

	return normalizedInput{
		CreateAsesoriaInput:        in,
		AcceptsDataProcessingBool:  acceptsDataProcessing,
		AuthorizesNotificationBool: authorizesNotification,
		HeadOfHouseholdBool:        headOfHousehold,
		ConsultationDateAt:         consultationDateAt,
		SubmittedAtTime:            submittedAtTime,
	}, nil
}

func parseOptionalYesNo(value string) (bool, error) {
	normalized := normalizeText(value)
	switch normalized {
	case "yes", "si", "s", "true", "1":
		return true, nil
	case "no", "n", "false", "0":
		return false, nil
	default:
		return false, fmt.Errorf("invalid yes/no value: %q", value)
	}
}

func parseDateTime(value string) (time.Time, error) {
	layouts := []string{
		time.RFC3339Nano,
		time.RFC3339,
		"2006-01-02 15:04:05",
		"2006-01-02T15:04",
		"2006-01-02",
	}

	for _, layout := range layouts {
		if ts, err := time.Parse(layout, value); err == nil {
			return ts, nil
		}
	}

	return time.Time{}, fmt.Errorf("invalid datetime %q", value)
}

func parseNonNegativeInt(value string) (int, error) {
	out, err := strconv.Atoi(strings.TrimSpace(value))
	if err != nil {
		return 0, err
	}
	if out < 0 {
		return 0, fmt.Errorf("value must be >= 0")
	}
	return out, nil
}

func normalizeText(value string) string {
	value = strings.TrimSpace(strings.ToLower(value))
	replacer := strings.NewReplacer(
		"\u00e1", "a",
		"\u00e9", "e",
		"\u00ed", "i",
		"\u00f3", "o",
		"\u00fa", "u",
	)
	return replacer.Replace(value)
}
