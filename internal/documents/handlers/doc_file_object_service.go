package handlers

import (
	"errors"
	"log"
	"strconv"

	"github.com/DuvanRozoParra/sicou/internal/documents/dto"
	services "github.com/DuvanRozoParra/sicou/internal/documents/service"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type SDocFileObjectHandler struct {
	service *services.SDocFileObjectService
}

func NewDocFileObjectHandler(service *services.SDocFileObjectService) *SDocFileObjectHandler {
	return &SDocFileObjectHandler{service: service}
}

func (h *SDocFileObjectHandler) UploadDoc(c *fiber.Ctx) error {

	// 🔹 1. Archivo
	fileHeader, err := c.FormFile("file")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "file is required")
	}

	file, err := fileHeader.Open()
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "cannot open file")
	}
	defer file.Close()

	// 🔹 2. document_kind_id (obligatorio)
	documentKindStr := c.FormValue("document_kind_id")
	if documentKindStr == "" {
		return fiber.NewError(fiber.StatusBadRequest, "document_kind_id is required")
	}

	documentKindInt, err := strconv.Atoi(documentKindStr)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "document_kind_id must be numeric")
	}

	// 🔹 3. Opcionales
	preturnoID := optionalString(c.FormValue("preturno_id"))
	caseID := optionalString(c.FormValue("case_id"))
	caseEventID := optionalString(c.FormValue("case_event_id"))
	notes := optionalString(c.FormValue("notes"))
	uploadedBy := optionalString(c.FormValue("uploaded_by"))

	// 🔹 Generar storage key
	storageKey := generateStorageKey(caseID, fileHeader.Filename)

	// 🔹 Construir request completo
	req := dto.CreateDocFileObjectRequest{
		DocumentKindID: int16(documentKindInt),
		PreturnoID:     preturnoID,
		CaseID:         caseID,
		CaseEventID:    caseEventID,
		Notes:          notes,
		UploadedBy:     uploadedBy,

		StorageKey:   storageKey,
		OriginalName: fileHeader.Filename,
		MimeType:     fileHeader.Header.Get("Content-Type"),
		SizeBytes:    fileHeader.Size,
	}

	doc, err := h.service.UploadDoc(c.Context(),
		file,
		fileHeader.Size,
		fileHeader.Filename,
		fileHeader.Header.Get("Content-Type"),
		req)
	if err != nil {

		switch {
		case errors.Is(err, services.ErrDuplicateDocument):
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error": "Ya existe un documento de este tipo para este expediente",
			})

		case errors.Is(err, services.ErrValidation):
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})

		default:
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Error interno al subir documento",
			})
		}
	}

	return c.Status(fiber.StatusCreated).JSON(doc)
}

func (h *SDocFileObjectHandler) UploadDocFolder(c *fiber.Ctx) error {

	// 1️⃣ Archivo
	fileHeader, err := c.FormFile("file")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "file is required")
	}

	file, err := fileHeader.Open()
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "cannot open file")
	}
	defer file.Close()

	// 2️⃣ document_kind_id (UUID)
	documentKindStr := c.FormValue("document_kind_id")
	if documentKindStr == "" {
		return fiber.NewError(fiber.StatusBadRequest, "document_kind_id is required")
	}

	documentKindInt64, err := strconv.ParseInt(documentKindStr, 10, 16)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "document_kind_id must be a valid number")
	}

	documentKindID := int16(documentKindInt64)

	if documentKindID <= 0 {
		return fiber.NewError(fiber.StatusBadRequest, "document_kind_id must be greater than 0")
	}

	// 3️⃣ folder_id (obligatorio)
	folderStr := c.FormValue("folder_id")
	if folderStr == "" {
		return fiber.NewError(fiber.StatusBadRequest, "folder_id is required")
	}

	folderID, err := uuid.Parse(folderStr)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "folder_id must be a valid UUID")
	}

	// 4️⃣ uploaded_by (UUID obligatorio)
	uploadedByStr := c.FormValue("uploaded_by")
	if uploadedByStr == "" {
		return fiber.NewError(fiber.StatusBadRequest, "uploaded_by is required")
	}

	uploadedBy, err := uuid.Parse(uploadedByStr)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "uploaded_by must be a valid UUID")
	}

	// 5️⃣ Opcionales
	var preturnoID *uuid.UUID
	if v := c.FormValue("preturno_id"); v != "" {
		parsed, err := uuid.Parse(v)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "preturno_id must be UUID")
		}
		preturnoID = &parsed
	}

	var caseID *uuid.UUID
	if v := c.FormValue("case_id"); v != "" {
		parsed, err := uuid.Parse(v)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "case_id must be UUID")
		}
		caseID = &parsed
	}

	var caseEventID *uuid.UUID
	if v := c.FormValue("case_event_id"); v != "" {
		parsed, err := uuid.Parse(v)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "case_event_id must be UUID")
		}
		caseEventID = &parsed
	}

	var notes *string
	if v := c.FormValue("notes"); v != "" {
		notes = &v
	}

	// 6️⃣ Construir request
	req := dto.CreateDocFileObjectFolderRequest{
		DocumentKindID: documentKindID,
		PreturnoID:     preturnoID,
		CaseID:         caseID,
		CaseEventID:    caseEventID,
		Notes:          notes,
		UploadedBy:     uploadedBy,
		FolderID:       folderID,
	}

	// 7️⃣ Llamar service
	doc, err := h.service.UploadDocFolder(
		c.Context(),
		file,
		fileHeader.Size,
		fileHeader.Filename,
		fileHeader.Header.Get("Content-Type"),
		req,
	)
	if err != nil {

		switch {
		case errors.Is(err, services.ErrDuplicateDocument):
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error": "Ya existe un documento de este tipo para este expediente",
			})

		case errors.Is(err, services.ErrValidation):
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})

		default:
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Error interno al subir documento",
			})
		}
	}

	return c.Status(fiber.StatusCreated).JSON(doc)
}

func (h *SDocFileObjectHandler) GetDocumentKinds(c *fiber.Ctx) error {

	log.Println("GetDocumentKinds called")

	kinds, err := h.service.GetDocumentKinds(c.Context())
	if err != nil {
		log.Printf("ERROR GetDocumentKinds: %v\n", err)

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "failed to fetch document kinds",
			"details": err.Error(), // quitar en producción si no quieres exponer detalles
		})
	}

	log.Printf("GetDocumentKinds success - found %d kinds\n", len(kinds))

	return c.Status(fiber.StatusOK).JSON(kinds)
}

func (h *SDocFileObjectHandler) GetCaseTimeLineItem(c *fiber.Ctx) error {

	// 🔹 1. Obtener case_id desde params
	caseIDParam := c.Params("case_id")

	if caseIDParam == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "case_id is required",
		})
	}

	// 🔹 2. Validar UUID
	caseID, err := uuid.Parse(caseIDParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid case_id format",
		})
	}

	// 🔹 3. Llamar service
	items, err := h.service.GetCaseTimeLineItem(c.Context(), caseID)
	if err != nil {

		switch err.Error() {

		case "case_id is required":
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})

		case "case not found":
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": err.Error(),
			})

		default:
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "internal server error",
			})
		}
	}

	// 🔹 4. Respuesta exitosa
	return c.Status(fiber.StatusOK).JSON(items)
}

func (h *SDocFileObjectHandler) ListCaseFiles(c *fiber.Ctx) error {

	caseIDParam := c.Params("case_id")
	if caseIDParam == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "case_id is required",
		})
	}

	caseID, err := uuid.Parse(caseIDParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid case_id format",
		})
	}

	files, err := h.service.ListCaseFiles(c.Context(), caseID)
	if err != nil {

		switch err.Error() {
		case "case_id is required":
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		case "case not found":
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "internal server error",
			})
		}
	}

	return c.Status(fiber.StatusOK).JSON(files)
}

func (h *SDocFileObjectHandler) ViewFile(c *fiber.Ctx) error {

	fileIDParam := c.Params("file_id")
	if fileIDParam == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "file_id is required",
		})
	}

	fileID, err := uuid.Parse(fileIDParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid file_id",
		})
	}

	meta, object, err := h.service.GetFileStream(c.Context(), fileID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "file not found",
		})
	}

	contentType := "application/octet-stream"
	if meta.MimeType != nil {
		contentType = *meta.MimeType
	}

	c.Set("Content-Type", contentType)
	c.Set("Content-Disposition", "inline; filename=\""+meta.OriginalName+"\"")

	return c.SendStream(object)
}

func (h *SDocFileObjectHandler) UpdateDocument(c *fiber.Ctx) error {

	log.Printf("Body raw: %s", c.Body())

	documentIDParam := c.Params("document_id")
	if documentIDParam == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "document_id is required",
		})
	}

	documentID, err := uuid.Parse(documentIDParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid document_id format",
		})
	}

	var req dto.UpdateDocumentRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}
	log.Printf("Parsed req: %+v", req)

	log.Printf("UpdateDocument - ID: %s", documentID.String())

	err = h.service.UpdateDocument(c.Context(), documentID, req)
	if err != nil {

		switch err.Error() {

		case "document_id is required":
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})

		case "at least one field must be provided":
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})

		case "original_name cannot be empty":
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})

		case "document not found":
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": err.Error(),
			})

		default:
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "internal server error",
			})
		}
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func (h *SDocFileObjectHandler) DeleteDocument(c *fiber.Ctx) error {

	documentIDParam := c.Params("document_id")
	if documentIDParam == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "document_id is required",
		})
	}

	documentID, err := uuid.Parse(documentIDParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid document_id format",
		})
	}

	err = h.service.DeleteDocument(c.Context(), documentID)
	if err != nil {

		switch err.Error() {

		case "document_id is required":
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})

		case "document not found":
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": err.Error(),
			})

		default:
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "internal server error",
			})
		}
	}

	return c.SendStatus(fiber.StatusNoContent)
}


func optionalString(value string) *string {
	if value == "" {
		return nil
	}
	return &value
}

func generateStorageKey(caseID *string, filename string) string {
	id := uuid.New().String()

	if caseID != nil {
		return "cases/" + *caseID + "/" + id + "-" + filename
	}

	return "misc/" + id + "-" + filename
}
