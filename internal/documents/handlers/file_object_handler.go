package handlers

import (
	"strconv"
	"strings"

	services "github.com/DuvanRozoParra/sicou/internal/documents/service"
	"github.com/gofiber/fiber/v2"
)

type FileObjectHandler struct {
	service *services.FileObjectService
}

func NewFileObjectHandler(service *services.FileObjectService) *FileObjectHandler {
	return &FileObjectHandler{service: service}
}

func (h *FileObjectHandler) Create(c *fiber.Ctx) error {

	// 1️⃣ Obtener archivo desde form-data
	fileHeader, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "file is required",
		})
	}

	file, err := fileHeader.Open()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to open file",
		})
	}
	defer file.Close()

	// 2️⃣ Llamar al service
	result, err := h.service.Create(
		c.Context(),
		file,
		fileHeader.Size,
		fileHeader.Filename,
		fileHeader.Header.Get("Content-Type"),
	)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(result)
}

func (h *FileObjectHandler) GetById(c *fiber.Ctx) error {
	id := c.Params("id")

	result, err := h.service.GetById(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "file not found",
		})
	}

	return c.JSON(result)
}

func (h *FileObjectHandler) GetAll(c *fiber.Ctx) error {
	result, err := h.service.GetAll(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "internal error",
		})
	}

	return c.JSON(result)
}

/*
func (h *FileObjectHandler) Update(c *fiber.Ctx) error {
	id := c.Params("id")

	var req dto.UpdateFileObjectRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	result, err := h.service.Update(c.Context(), id, req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "update failed",
		})
	}

	return c.JSON(result)
}
*/

func (h *FileObjectHandler) Delete(c *fiber.Ctx) error {
	id := c.Params("id")

	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "id is required",
		})
	}

	err := h.service.Delete(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func (h *FileObjectHandler) Download(c *fiber.Ctx) error {
	id := c.Params("id")

	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "id is required",
		})
	}

	object, metadata, err := h.service.GetFileStream(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "file not found",
		})
	}

	contentType := strings.TrimSpace(metadata.MimeType)
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	filename := sanitizeHeaderFilename(metadata.OriginalName)
	c.Set("Content-Type", contentType)
	c.Set("Content-Disposition", "attachment; filename=\""+filename+"\"")

	return c.SendStream(object)
}

func (h *FileObjectHandler) View(c *fiber.Ctx) error {
	id := c.Params("id")

	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "id is required",
		})
	}

	object, metadata, err := h.service.GetFileStream(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "file not found",
		})
	}

	contentType := strings.TrimSpace(metadata.MimeType)
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	filename := sanitizeHeaderFilename(metadata.OriginalName)
	c.Set("Content-Type", contentType)
	c.Set("Content-Disposition", "inline; filename=\""+filename+"\"")
	c.Set("Accept-Ranges", "bytes")

	return c.SendStream(object)
}

func (h *FileObjectHandler) GetPresignedURL(c *fiber.Ctx) error {
	id := c.Params("id")

	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "id is required",
		})
	}

	url, err := h.service.GetPresignedURL(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "file not found",
		})
	}

	return c.JSON(fiber.Map{
		"url": url,
	})
}

func (h *FileObjectHandler) List(c *fiber.Ctx) error {

	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	files, err := h.service.ListDocuments(c.Context(), page, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to list documents",
		})
	}

	return c.JSON(files)
}

func sanitizeHeaderFilename(value string) string {
	name := strings.TrimSpace(value)
	if name == "" {
		return "file.bin"
	}
	name = strings.ReplaceAll(name, "\\", "\\\\")
	name = strings.ReplaceAll(name, "\"", "\\\"")
	return name
}
