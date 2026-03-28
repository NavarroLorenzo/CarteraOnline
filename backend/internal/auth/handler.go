package auth

import "github.com/gin-gonic/gin"

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Register(c *gin.Context) {
	var input RegisterInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{
			"error":   "Datos inválidos",
			"details": err.Error(),
		})
		return
	}

	response, err := h.service.Register(input)
	if err != nil {
		switch err {
		case ErrEmailAlreadyInUse, ErrUsernameAlreadyInUse:
			c.JSON(409, gin.H{"error": err.Error()})
		default:
			c.JSON(500, gin.H{"error": "No se pudo registrar el usuario"})
		}
		return
	}

	c.JSON(201, response)
}

func (h *Handler) Login(c *gin.Context) {
	var input LoginInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{
			"error":   "Datos inválidos",
			"details": err.Error(),
		})
		return
	}

	response, err := h.service.Login(input)
	if err != nil {
		switch err {
		case ErrInvalidCredentials:
			c.JSON(401, gin.H{"error": err.Error()})
		default:
			c.JSON(500, gin.H{"error": "No se pudo iniciar sesión"})
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
			c.JSON(401, gin.H{"error": "Usuario no encontrado"})
		default:
			c.JSON(500, gin.H{"error": "No se pudo obtener el usuario autenticado"})
		}
		return
	}

	c.JSON(200, user)
}
