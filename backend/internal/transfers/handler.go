package transfers

import (
	"cartera-app/backend/internal/auth"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Create(c *gin.Context) {
	userID, ok := auth.AbortIfUnauthenticated(c)
	if !ok {
		return
	}

	var input CreateTransferInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{
			"error":   "Datos inválidos",
			"details": err.Error(),
		})
		return
	}

	err := h.service.CreateTransfer(userID, input)
	if err != nil {
		switch err {
		case ErrSameAccount:
			c.JSON(400, gin.H{"error": err.Error()})
		case ErrAccountNotFound:
			c.JSON(400, gin.H{"error": err.Error()})
		default:
			c.JSON(500, gin.H{"error": "No se pudo realizar la transferencia"})
		}
		return
	}

	c.JSON(201, gin.H{
		"message": "Transferencia realizada correctamente",
	})
}
