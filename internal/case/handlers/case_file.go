package handlers

import (
	"log"
	"strconv"
	"strings"

	"github.com/DuvanRozoParra/sicou/internal/case/dto"
	"github.com/DuvanRozoParra/sicou/internal/case/service"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
)

type SCaseFileHandler struct {
	service *service.SCaseFileService
}

func NewCaseFileHandler(service *service.SCaseFileService) *SCaseFileHandler {
	return &SCaseFileHandler{service: service}
}

func (s *SCaseFileHandler) Create(c *fiber.Ctx) error {
	var req dto.CreateCaseFileRequest

	if err := c.BodyParser(&req); err != nil {
		log.Println("BodyParser error:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	result, err := s.service.Create(c.Context(), req)
	if err != nil {

		if pgErr, ok := err.(*pgconn.PgError); ok {
			if pgErr.Code == "23505" && pgErr.ConstraintName == "case_file_preturno_id_key" {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"error": "no se puede tener más de un caso en un turno",
				})
			}
		}

		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.Status(fiber.StatusCreated).JSON(result)
}

func (s *SCaseFileHandler) Update(c *fiber.Ctx) error {
	idParam := c.Params("id")

	id, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid id",
		})
	}

	var req dto.UpdateCaseFileRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	result, err := s.service.Update(c.Context(), id, req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(result)
}
func (s *SCaseFileHandler) List(c *fiber.Ctx) error {

	limitStr := c.Query("limit", "10")
	offsetStr := c.Query("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 10
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		offset = 0
	}

	result, err := s.service.List(c.Context(), limit, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(result)
}
func (s *SCaseFileHandler) GetById(c *fiber.Ctx) error {
	idParam := c.Params("id")

	id, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid id",
		})
	}

	result, err := s.service.GetById(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "case file not found",
		})
	}

	return c.JSON(result)
}
func (s *SCaseFileHandler) Delete(c *fiber.Ctx) error {
	idParam := c.Params("id")

	id, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid id",
		})
	}

	result, err := s.service.Delete(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "case file not found",
		})
	}

	return c.JSON(result)
}
func (s *SCaseFileHandler) ListCase(c *fiber.Ctx) error {

	result, err := s.service.ListCase(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(result)
}

func (h *SCaseFileHandler) ListCasesByEmail(c *fiber.Ctx) error {

	email := strings.TrimSpace(c.Query("email"))
	if email == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "email es requerido",
		})
	}

	cases, err := h.service.ListCasesByEmail(
		c.Context(),
		email,
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": cases,
	})
}

func (h *SCaseFileHandler) CaseOTP(c *fiber.Ctx) error {

	var req dto.CaseOtpRequest

	// 1️⃣ Parsear body
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Solicitud inválida",
		})
	}

	log.Printf("Email recibido: %+v\n", req.Email)

	// 2️⃣ Validación mínima defensiva
	if req.Email == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "El correo es requerido",
		})
	}

	// 3️⃣ Llamar service
	resp, err := h.service.CaseOTP(c.Context(), req)
	if err != nil {

		// Puedes mejorar esto luego con manejo centralizado de errores
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	// 4️⃣ Respuesta exitosa
	return c.Status(fiber.StatusOK).JSON(resp)
}

func (h *SCaseFileHandler) VerifyOTP(c *fiber.Ctx) error {

	var req dto.VerifyCaseFileOtpRequest

	// 1️⃣ Parsear body
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "payload inválido",
		})
	}

	// 2️⃣ Validación básica (defensiva)
	if strings.TrimSpace(req.Email) == "" ||
		strings.TrimSpace(req.OTP) == "" {

		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "citizen_id y otp son requeridos",
		})
	}

	// 3️⃣ Llamar al service
	err := h.service.VerifyOtp(
		c.Context(),
		req.Email,
		req.OTP,
	)

	if err != nil {

		// Puedes mapear errores específicos aquí
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	// 4️⃣ Respuesta exitosa
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "OTP verificado correctamente",
	})
}

func (h *SCaseFileHandler) FolderList(c *fiber.Ctx) error {
	ctx := c.UserContext()

	folders, err := h.service.FolderList(ctx)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to retrieve folders",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": folders,
	})
}
