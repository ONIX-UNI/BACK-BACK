package service

import (
	"fmt"
	stdmail "net/mail"
	"strings"
	"time"
)

type normalizedInput struct {
	CreatePQRSInput
	QueryTypeIsAnonymous         bool
	AllowsElectronicResponseBool bool
	RequestAgainstStudentBool    bool
	AcceptsDataProcessingBool    bool
	Email                        string
	RequestDescription           string
	FirstName                    string
	FirstLastName                string
	SubmittedAt                  time.Time
}

func normalizeAndValidate(in CreatePQRSInput) (normalizedInput, *ValidationError) {
	fieldsErr := map[string]string{}

	in.QueryType = strings.TrimSpace(in.QueryType)
	in.PersonType = strings.TrimSpace(in.PersonType)
	in.DocumentType = strings.TrimSpace(in.DocumentType)
	in.DocumentOrTaxID = strings.TrimSpace(in.DocumentOrTaxID)
	in.FirstName = strings.TrimSpace(in.FirstName)
	in.MiddleName = strings.TrimSpace(in.MiddleName)
	in.FirstLastName = strings.TrimSpace(in.FirstLastName)
	in.SecondLastName = strings.TrimSpace(in.SecondLastName)
	in.Gender = strings.TrimSpace(in.Gender)
	in.Address = strings.TrimSpace(in.Address)
	in.NeighborhoodArea = strings.TrimSpace(in.NeighborhoodArea)
	in.AllowsElectronicResponse = strings.TrimSpace(in.AllowsElectronicResponse)
	in.Email = strings.TrimSpace(in.Email)
	in.Phone = strings.TrimSpace(in.Phone)
	in.PopulationGroup = strings.TrimSpace(in.PopulationGroup)
	in.OtherPopulationGroup = strings.TrimSpace(in.OtherPopulationGroup)
	in.RequestDescription = strings.TrimSpace(in.RequestDescription)
	in.ResponseChannel = strings.TrimSpace(in.ResponseChannel)
	in.RequestType = strings.TrimSpace(in.RequestType)
	in.RequestAgainstStudent = strings.TrimSpace(in.RequestAgainstStudent)
	in.ResponsibleStudentName = strings.TrimSpace(in.ResponsibleStudentName)
	in.ResponsibleStudentProgram = strings.TrimSpace(in.ResponsibleStudentProgram)
	in.StudentCaseDescription = strings.TrimSpace(in.StudentCaseDescription)
	in.AcceptsDataProcessing = strings.TrimSpace(in.AcceptsDataProcessing)
	in.SubmittedAt = strings.TrimSpace(in.SubmittedAt)

	queryTypeIsAnonymous, errQueryType := parseQueryType(in.QueryType)
	if errQueryType != nil {
		fieldsErr["queryType"] = "must indicate if query is anonymous or identified"
	}

	allowsElectronicResponse, errAllowsElectronicResponse := parseYesNo(in.AllowsElectronicResponse)
	if errAllowsElectronicResponse != nil {
		fieldsErr["allowsElectronicResponse"] = "must be Yes or No"
	}

	requestAgainstStudent, errRequestAgainstStudent := parseYesNo(in.RequestAgainstStudent)
	if errRequestAgainstStudent != nil {
		fieldsErr["requestAgainstStudent"] = "must be Yes or No"
	}

	acceptsDataProcessing, errAcceptsDataProcessing := parseYesNo(in.AcceptsDataProcessing)
	if errAcceptsDataProcessing != nil {
		fieldsErr["acceptsDataProcessing"] = "must be Yes or No"
	}
	if !acceptsDataProcessing {
		fieldsErr["acceptsDataProcessing"] = "must accept data processing"
	}

	if in.RequestDescription == "" {
		fieldsErr["requestDescription"] = "is required"
	}
	if in.ResponseChannel == "" {
		fieldsErr["responseChannel"] = "is required"
	}
	if in.RequestType == "" {
		fieldsErr["requestType"] = "is required"
	}
	if in.SubmittedAt == "" {
		fieldsErr["submittedAt"] = "is required"
	}

	submittedAt, errSubmittedAt := parseSubmittedAt(in.SubmittedAt)
	if errSubmittedAt != nil {
		fieldsErr["submittedAt"] = "invalid date format"
	}

	if in.Email != "" {
		if _, err := stdmail.ParseAddress(in.Email); err != nil {
			fieldsErr["email"] = "invalid email"
		}
	}

	if errQueryType == nil && !queryTypeIsAnonymous {
		requireIfEmpty(fieldsErr, "personType", in.PersonType)
		requireIfEmpty(fieldsErr, "documentType", in.DocumentType)
		requireIfEmpty(fieldsErr, "documentOrTaxId", in.DocumentOrTaxID)
		requireIfEmpty(fieldsErr, "firstName", in.FirstName)
		requireIfEmpty(fieldsErr, "firstLastName", in.FirstLastName)
		requireIfEmpty(fieldsErr, "gender", in.Gender)
		requireIfEmpty(fieldsErr, "address", in.Address)
		requireIfEmpty(fieldsErr, "neighborhoodArea", in.NeighborhoodArea)
		requireIfEmpty(fieldsErr, "email", in.Email)
		requireIfEmpty(fieldsErr, "phone", in.Phone)
	}

	if allowsElectronicResponse && in.Email == "" {
		fieldsErr["email"] = "is required when allowsElectronicResponse is Yes"
	}

	populationGroupNormalized := normalizeText(in.PopulationGroup)
	if populationGroupNormalized == "other" || populationGroupNormalized == "otro" {
		requireIfEmpty(fieldsErr, "otherPopulationGroup", in.OtherPopulationGroup)
	}

	if requestAgainstStudent {
		requireIfEmpty(fieldsErr, "responsibleStudentName", in.ResponsibleStudentName)
		requireIfEmpty(fieldsErr, "responsibleStudentProgram", in.ResponsibleStudentProgram)
		requireIfEmpty(fieldsErr, "studentCaseDescription", in.StudentCaseDescription)
	}

	if len(fieldsErr) > 0 {
		return normalizedInput{}, &ValidationError{Fields: fieldsErr}
	}

	return normalizedInput{
		CreatePQRSInput:              in,
		QueryTypeIsAnonymous:         queryTypeIsAnonymous,
		AllowsElectronicResponseBool: allowsElectronicResponse,
		RequestAgainstStudentBool:    requestAgainstStudent,
		AcceptsDataProcessingBool:    acceptsDataProcessing,
		Email:                        in.Email,
		RequestDescription:           in.RequestDescription,
		FirstName:                    in.FirstName,
		FirstLastName:                in.FirstLastName,
		SubmittedAt:                  submittedAt,
	}, nil
}

func requireIfEmpty(fieldsErr map[string]string, fieldName string, value string) {
	if strings.TrimSpace(value) == "" {
		fieldsErr[fieldName] = "is required"
	}
}

func parseQueryType(value string) (bool, error) {
	normalized := normalizeText(value)
	switch normalized {
	case "anonymous", "anonima", "anonimo", "anon":
		return true, nil
	case "identified", "identificada", "identificado":
		return false, nil
	case "yes", "si", "s", "true", "1":
		return true, nil
	case "no", "n", "false", "0":
		return false, nil
	default:
		return false, fmt.Errorf("invalid queryType value: %q", value)
	}
}

func parseYesNo(value string) (bool, error) {
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

func parseSubmittedAt(value string) (time.Time, error) {
	layouts := []string{
		time.RFC3339Nano,
		time.RFC3339,
		"2006-01-02 15:04:05",
		"2006-01-02",
	}

	for _, layout := range layouts {
		if ts, err := time.Parse(layout, value); err == nil {
			return ts, nil
		}
	}

	return time.Time{}, fmt.Errorf("invalid submittedAt %q", value)
}
