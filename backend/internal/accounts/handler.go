package accounts

import "github.com/gin-gonic/gin"

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{
		service: service,
	}
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

	account := h.service.Create(input)
	c.JSON(201, account)
}

func (h *Handler) GetAll(c *gin.Context) {
	accounts := h.service.GetAll()
	c.JSON(200, accounts)
}
