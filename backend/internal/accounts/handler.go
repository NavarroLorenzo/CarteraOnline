package accounts

import "github.com/gin-gonic/gin"

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Create(c *gin.Context) {
	var input CreateAccountInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{
			"error":   "Datos inválidos",
			"details": err.Error(),
		})
		return
	}

	account, err := h.service.Create(input)
	if err != nil {
		c.JSON(500, gin.H{"error": "No se pudo crear la cuenta"})
		return
	}

	c.JSON(201, account)
}

func (h *Handler) GetAll(c *gin.Context) {
	accounts, err := h.service.GetAll()
	if err != nil {
		c.JSON(500, gin.H{"error": "No se pudieron obtener las cuentas"})
		return
	}

	c.JSON(200, accounts)
}
