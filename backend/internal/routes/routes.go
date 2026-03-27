package routes

import (
	"cartera-app/backend/internal/accounts"
	"cartera-app/backend/internal/transactions"
	"cartera-app/backend/internal/transfers"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func SetupRouter(db *pgxpool.Pool) *gin.Engine {
	router := gin.Default()

	api := router.Group("/api")
	{
		accountRepo := accounts.NewPostgresRepository(db)
		transactionRepo := transactions.NewPostgresRepository(db)

		accountService := accounts.NewService(accountRepo)
		transactionService := transactions.NewService(transactionRepo, accountService, accountService)

		accountHandler := accounts.NewHandler(accountService, transactionService)
		transactionHandler := transactions.NewHandler(transactionService)

		transferService := transfers.NewService(accountService, transactionService)
		transferHandler := transfers.NewHandler(transferService)

		api.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "API funcionando correctamente"})
		})

		api.POST("/accounts", accountHandler.Create)
		api.GET("/accounts", accountHandler.GetAll)

		api.POST("/transactions", transactionHandler.Create)
		api.GET("/transactions", transactionHandler.GetAll)
		api.GET("/transactions/balance", transactionHandler.GetBalance)
		api.GET("/transactions/balance-by-account", transactionHandler.GetBalanceByAccount)
		api.GET("/transactions/summary", transactionHandler.GetSummary)
		api.GET("/transactions/:id", transactionHandler.GetByID)
		api.PUT("/transactions/:id", transactionHandler.Update)
		api.DELETE("/transactions/:id", transactionHandler.Delete)

		api.POST("/transfers", transferHandler.Create)
	}

	return router
}
