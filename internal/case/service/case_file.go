package service

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"regexp"
	"strings"
	"time"

	"github.com/DuvanRozoParra/sicou/internal/case/dto"
	"github.com/DuvanRozoParra/sicou/internal/case/models"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type SCaseFileService struct {
	repo models.ICaseFile
}

func NewCaseFileService(repo models.ICaseFile) *SCaseFileService {
	return &SCaseFileService{repo: repo}
}

func (s *SCaseFileService) Create(ctx context.Context, req dto.CreateCaseFileRequest) (*dto.CaseFile, error) {
	if req.CitizenID == uuid.Nil {
		return nil, errors.New("citizen_id is required")
	}

	if req.PreturnoID == uuid.Nil {
		return nil, errors.New("preturno_id is required")
	}

	if req.ServiceTypeID == 0 {
		return nil, errors.New("service_type_id is required")
	}

	return s.repo.Create(ctx, req)
}
func (s *SCaseFileService) Update(ctx context.Context, id uuid.UUID, req dto.UpdateCaseFileRequest) (*dto.CaseFile, error) {
	if id == uuid.Nil {
		return nil, errors.New("invalid id")
	}

	// Validación opcional de estado permitido
	if req.Status != nil {
		switch *req.Status {
		case "ABIERTO", "PENDIENTE_DOCUMENTOS", "EN_TRAMITE", "DESCARGADO", "CERRADO":
			// válido
		default:
			return nil, errors.New("invalid status value")
		}
	}

	return s.repo.Update(ctx, id, req)
}
func (s *SCaseFileService) List(ctx context.Context, limit, offset int) ([]dto.CaseFile, error) {
	if limit <= 0 {
		limit = 10
	}

	if offset < 0 {
		offset = 0
	}

	return s.repo.List(ctx, limit, offset)
}
func (s *SCaseFileService) GetById(ctx context.Context, id uuid.UUID) (*dto.CaseFile, error) {
	if id == uuid.Nil {
		return nil, errors.New("invalid id")
	}

	return s.repo.GetById(ctx, id)
}
func (s *SCaseFileService) Delete(ctx context.Context, id uuid.UUID) (*dto.CaseFile, error) {

	if id == uuid.Nil {
		return nil, errors.New("invalid id")
	}

	return s.repo.Delete(ctx, id)
}

func (s *SCaseFileService) ListCase(ctx context.Context) ([]dto.CaseListItem, error) {
	return s.repo.ListCases(ctx)
}

func (s *SCaseFileService) ListCasesByEmail(
	ctx context.Context,
	email string,
) ([]dto.CaseListItem, error) {

	// 1️⃣ Validaciones básicas
	email = strings.TrimSpace(email)
	if email == "" {
		return nil, errors.New("email es requerido")
	}

	// validación mínima formato
	if !isValidEmail(email) {
		return nil, errors.New("email inválido")
	}

	// 2️⃣ Delegar al repository
	cases, err := s.repo.ListCasesByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	return cases, nil
}

func (s *SCaseFileService) CaseOTP(
	ctx context.Context,
	req dto.CaseOtpRequest,
) (*dto.CaseOtpResponse, error) {

	// 1️⃣ Validación correcta para string
	if req.Email == "" {
		return nil, errors.New("el correo es requerido")
	}

	// 2️⃣ Validar existencia
	citizenID, err := s.repo.GetCitizenIDByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}

	// 3️⃣ Generar OTP
	otp, err := GenerateNumericOTP(6)
	if err != nil {
		return nil, err
	}

	// 4️⃣ Hash
	hashedOtp, err := bcrypt.GenerateFromPassword(
		[]byte(otp),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return nil, err
	}

	// 5️⃣ Expiración
	expiresAt := time.Now().Add(10 * time.Minute)

	// 6️⃣ Guardar
	err = s.repo.SaveOtp(
		ctx,
		citizenID,
		string(hashedOtp),
		expiresAt,
	)
	if err != nil {
		return nil, err
	}

	// 📩 Crear email en outbox
	subject := "Código de verificación"

	body := fmt.Sprintf(`
Hola,

Tu código de verificación es:

%s

Este código expira en 10 minutos.

Si no solicitaste este código, puedes ignorar este mensaje.
`, otp)

	err = s.repo.InsertEmailOutbox(
		ctx,
		[]string{req.Email},
		subject,
		body,
	)
	if err != nil {
		return nil, err
	}

	return &dto.CaseOtpResponse{
		Message: "Código enviado al correo.",
	}, nil
}

func (s *SCaseFileService) VerifyOtp(
	ctx context.Context,
	email string,
	otp string,
) error {
	// 1️⃣ Validaciones básicas
	if strings.TrimSpace(email) == "" {
		return errors.New("citizen_id es requerido")
	}

	if strings.TrimSpace(otp) == "" {
		return errors.New("otp es requerido")
	}

	// 2️⃣ Validar formato OTP
	if len(otp) != 6 {
		return errors.New("otp debe tener 6 dígitos")
	}

	if !isNumeric(otp) {
		return errors.New("otp debe ser numérico")
	}

	citizenID, err := s.repo.GetCitizenIDByEmail(ctx, email)
	if err != nil {
		return err
	}

	// 3️⃣ Delegar al repository
	err = s.repo.VerifyOtp(ctx, citizenID, otp)
	if err != nil {
		return err
	}

	return nil
}

func (s *SCaseFileService) FolderList(ctx context.Context) ([]dto.FolderResponse, error) {
	return s.repo.FolderList(ctx)
}

func isNumeric(value string) bool {
	for _, r := range value {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}

func GenerateNumericOTP(length int) (string, error) {
	max := big.NewInt(10)

	otp := make([]byte, length)

	for i := 0; i < length; i++ {
		n, err := rand.Int(rand.Reader, max)
		if err != nil {
			return "", err
		}
		otp[i] = byte(n.Int64()) + '0'
	}

	return fmt.Sprintf("%s", otp), nil
}

func isValidEmail(email string) bool {
	// validación simple, suficiente para backend
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(email)
}
