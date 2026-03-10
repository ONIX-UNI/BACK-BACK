package repository

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type internalEmailPayload struct {
	QueryType                 string                    `json:"queryType"`
	PersonType                string                    `json:"personType"`
	DocumentType              string                    `json:"documentType"`
	DocumentOrTaxID           string                    `json:"documentOrTaxId"`
	FirstName                 string                    `json:"firstName"`
	MiddleName                string                    `json:"middleName"`
	FirstLastName             string                    `json:"firstLastName"`
	SecondLastName            string                    `json:"secondLastName"`
	Gender                    string                    `json:"gender"`
	Address                   string                    `json:"address"`
	NeighborhoodArea          string                    `json:"neighborhoodArea"`
	AllowsElectronicResponse  string                    `json:"allowsElectronicResponse"`
	Email                     string                    `json:"email"`
	Phone                     string                    `json:"phone"`
	PopulationGroup           string                    `json:"populationGroup"`
	OtherPopulationGroup      string                    `json:"otherPopulationGroup"`
	RequestDescription        string                    `json:"requestDescription"`
	ResponseChannel           string                    `json:"responseChannel"`
	RequestType               string                    `json:"requestType"`
	RequestAgainstStudent     string                    `json:"requestAgainstStudent"`
	ResponsibleStudentName    string                    `json:"responsibleStudentName"`
	ResponsibleStudentProgram string                    `json:"responsibleStudentProgram"`
	StudentCaseDescription    string                    `json:"studentCaseDescription"`
	AcceptsDataProcessing     string                    `json:"acceptsDataProcessing"`
	SubmittedAt               string                    `json:"submittedAt"`
	Attachments               []internalEmailAttachment `json:"attachments"`
}

type internalEmailAttachment struct {
	Name string `json:"name"`
	Type string `json:"type"`
	Size int64  `json:"size"`
	URL  string `json:"url"`
}

func buildInternalEmailBody(
	radicado string,
	subject string,
	citizenName string,
	fromEmail *string,
	bodyPayload string,
) string {
	email := ""
	if fromEmail != nil {
		email = strings.TrimSpace(*fromEmail)
	}

	payload, _ := parseInternalEmailPayload(bodyPayload)

	citizen := strings.TrimSpace(citizenName)
	payloadFullName := buildPayloadFullName(payload)
	if payloadFullName != "" {
		citizen = payloadFullName
	} else if citizen == "" || strings.EqualFold(citizen, "citizen") {
		citizen = "No informado"
	}
	if citizen == "" {
		citizen = "No informado"
	}

	if email == "" {
		email = strings.TrimSpace(payload.Email)
	}
	if email == "" {
		email = "No informado"
	}

	asunto := strings.TrimSpace(subject)
	if asunto == "" {
		asunto = strings.TrimSpace(payload.RequestDescription)
	}
	if asunto == "" {
		asunto = "Sin asunto"
	}

	requestType := strings.TrimSpace(payload.RequestType)
	submittedAt := formatSubmittedAt(payload.SubmittedAt)
	responseChannel := strings.TrimSpace(payload.ResponseChannel)
	documentLine := formatDocumentLine(payload.DocumentType, payload.DocumentOrTaxID)
	allowsElectronic := formatYesNoForEmail(payload.AllowsElectronicResponse)
	acceptsData := formatYesNoForEmail(payload.AcceptsDataProcessing)
	requestAgainstStudent := formatYesNoForEmail(payload.RequestAgainstStudent)
	description := strings.TrimSpace(payload.RequestDescription)
	if description == "" {
		description = asunto
	}

	var b strings.Builder
	b.WriteString("Se recibio una nueva PQRS en SICOU.\n\n")
	appendLabeledLine(&b, "Radicado", strings.TrimSpace(radicado))
	appendLabeledLine(&b, "Tipo de solicitud", requestType)
	appendLabeledLine(&b, "Asunto", asunto)
	appendLabeledLine(&b, "Fecha de envio", submittedAt)
	appendLabeledLine(&b, "Canal de respuesta preferido", responseChannel)

	b.WriteString("\nDatos del solicitante:\n")
	appendBulletLine(&b, "Nombre completo", citizen)
	appendBulletLine(&b, "Documento", documentLine)
	appendBulletLine(&b, "Tipo de persona", payload.PersonType)
	appendBulletLine(&b, "Genero", payload.Gender)
	appendBulletLine(&b, "Telefono", payload.Phone)
	appendBulletLine(&b, "Correo electronico", email)
	appendBulletLine(&b, "Direccion", payload.Address)
	appendBulletLine(&b, "Barrio o sector", payload.NeighborhoodArea)
	appendBulletLine(&b, "Grupo poblacional", payload.PopulationGroup)
	appendBulletLine(&b, "Otro grupo poblacional", payload.OtherPopulationGroup)
	appendBulletLine(&b, "Acepta respuesta electronica", allowsElectronic)
	appendBulletLine(&b, "Autoriza tratamiento de datos", acceptsData)

	b.WriteString("\nDescripcion de la solicitud:\n")
	b.WriteString(description)
	b.WriteString("\n")

	if requestAgainstStudent != "" || strings.TrimSpace(payload.ResponsibleStudentName) != "" ||
		strings.TrimSpace(payload.ResponsibleStudentProgram) != "" ||
		strings.TrimSpace(payload.StudentCaseDescription) != "" {
		b.WriteString("\nInformacion sobre posible estudiante involucrado:\n")
		appendBulletLine(&b, "Aplica", requestAgainstStudent)
		appendBulletLine(&b, "Nombre del estudiante", payload.ResponsibleStudentName)
		appendBulletLine(&b, "Programa del estudiante", payload.ResponsibleStudentProgram)
		appendBulletLine(&b, "Descripcion del caso", payload.StudentCaseDescription)
	}

	if len(payload.Attachments) > 0 {
		b.WriteString("\nAdjuntos:\n")
		for _, attachment := range payload.Attachments {
			name := firstNonEmpty(strings.TrimSpace(attachment.Name), "Adjunto sin nombre")
			b.WriteString("- ")
			b.WriteString(name)
			b.WriteString("\n")
		}
	} else {
		b.WriteString("\nAdjuntos: No se registraron adjuntos.\n")
	}

	b.WriteString("\nPor favor gestionar esta PQRS conforme al flujo del consultorio.\n")
	return b.String()
}

func buildCitizenAckBody(radicado string, citizenName string) string {
	name := strings.TrimSpace(citizenName)
	if name == "" {
		name = "ciudadano"
	}

	return fmt.Sprintf(
		"Hola %s,\n\nRecibimos tu PQRS correctamente.\nTu numero de radicado es %s.\n\nCon este radicado puedes hacer seguimiento al caso.\n",
		name,
		radicado,
	)
}

func formatRadicado(createdAt time.Time, sequence int64) string {
	year := createdAt.Year()
	return fmt.Sprintf("PQRS-%d-%04d", year, sequence)
}

func parseInternalEmailPayload(raw string) (internalEmailPayload, bool) {
	value := strings.TrimSpace(raw)
	if value == "" {
		return internalEmailPayload{}, false
	}

	var payload internalEmailPayload
	if err := json.Unmarshal([]byte(value), &payload); err != nil {
		return internalEmailPayload{}, false
	}

	return payload, true
}

func buildPayloadFullName(payload internalEmailPayload) string {
	parts := []string{
		strings.TrimSpace(payload.FirstName),
		strings.TrimSpace(payload.MiddleName),
		strings.TrimSpace(payload.FirstLastName),
		strings.TrimSpace(payload.SecondLastName),
	}

	nonEmpty := make([]string, 0, len(parts))
	for _, part := range parts {
		if part == "" {
			continue
		}
		nonEmpty = append(nonEmpty, part)
	}

	return strings.TrimSpace(strings.Join(nonEmpty, " "))
}

func appendLabeledLine(builder *strings.Builder, label string, value string) {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return
	}
	builder.WriteString(label)
	builder.WriteString(": ")
	builder.WriteString(trimmed)
	builder.WriteString("\n")
}

func appendBulletLine(builder *strings.Builder, label string, value string) {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return
	}
	builder.WriteString("- ")
	builder.WriteString(label)
	builder.WriteString(": ")
	builder.WriteString(trimmed)
	builder.WriteString("\n")
}

func formatDocumentLine(documentType string, documentNumber string) string {
	docType := strings.TrimSpace(documentType)
	docNumber := strings.TrimSpace(documentNumber)
	if docType == "" {
		return docNumber
	}
	if docNumber == "" {
		return docType
	}
	return fmt.Sprintf("%s %s", docType, docNumber)
}

func formatSubmittedAt(value string) string {
	ts := strings.TrimSpace(value)
	if ts == "" {
		return ""
	}

	layouts := []string{
		time.RFC3339Nano,
		time.RFC3339,
		"2006-01-02 15:04:05",
		"2006-01-02",
	}

	for _, layout := range layouts {
		parsed, err := time.Parse(layout, ts)
		if err != nil {
			continue
		}
		return parsed.In(time.Local).Format("02/01/2006 15:04")
	}

	return ts
}

func formatYesNoForEmail(value string) string {
	normalized := normalizeEmailText(value)
	switch normalized {
	case "yes", "si", "s", "true", "1":
		return "Si"
	case "no", "n", "false", "0":
		return "No"
	default:
		return strings.TrimSpace(value)
	}
}

func normalizeEmailText(value string) string {
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

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		trimmed := strings.TrimSpace(value)
		if trimmed == "" {
			continue
		}
		return trimmed
	}
	return ""
}
