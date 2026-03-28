package auth

import (
	"cartera-app/backend/internal/shared/httpjson"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Register(c *gin.Context) {
	var input RegisterInput

	if err := c.ShouldBindJSON(&input); err != nil {
		httpjson.ErrorWithDetails(c, 400, "invalid_request", "Datos inválidos", err.Error())
		return
	}

	response, err := h.service.Register(input)
	if err != nil {
		switch err {
		case ErrEmailAlreadyInUse, ErrUsernameAlreadyInUse:
			httpjson.Error(c, 409, "auth_conflict", err.Error())
		case ErrInvalidEmail:
			httpjson.Error(c, 400, "invalid_email", err.Error())
		case ErrInvalidUsername:
			httpjson.Error(c, 400, "invalid_username", err.Error())
		case ErrWeakPassword:
			httpjson.Error(c, 400, "weak_password", err.Error())
		default:
			httpjson.Error(c, 500, "register_failed", "No se pudo registrar el usuario")
		}
		return
	}

	c.JSON(201, response)
}

func (h *Handler) Login(c *gin.Context) {
	var input LoginInput

	if err := c.ShouldBindJSON(&input); err != nil {
		httpjson.ErrorWithDetails(c, 400, "invalid_request", "Datos inválidos", err.Error())
		return
	}

	response, err := h.service.Login(input)
	if err != nil {
		switch err {
		case ErrInvalidCredentials:
			httpjson.Error(c, 401, "invalid_credentials", err.Error())
		default:
			httpjson.Error(c, 500, "login_failed", "No se pudo iniciar sesión")
		}
		return
	}

	c.JSON(200, response)
}

func (h *Handler) Me(c *gin.Context) {
	userID, ok := AbortIfUnauthenticated(c)
	if !ok {
		return
	}

	user, err := h.service.GetByID(userID)
	if err != nil {
		switch err {
		case ErrUserNotFound:
			httpjson.Error(c, 401, "user_not_found", "Usuario no encontrado")
		default:
			httpjson.Error(c, 500, "auth_me_failed", "No se pudo obtener el usuario autenticado")
		}
		return
	}

	c.JSON(200, user)
}
