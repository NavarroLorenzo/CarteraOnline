package routes

import (
	"cartera-app/backend/internal/accounts"
	"cartera-app/backend/internal/transactions"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func SetupRouter(db *pgxpool.Pool) *gin.Engine {
	router := gin.Default()

	api := router.Group("/api")
	{
		accountRepo := accounts.NewPostgresRepository(db)
		accountService := accounts.NewService(accountRepo)
		accountHandler := accounts.NewHandler(accountService)

		transactionRepo := transactions.NewPostgresRepository(db)
		transactionService := transactions.NewService(transactionRepo, accountService, accountService)
		transactionHandler := transactions.NewHandler(transactionService)

		api.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "API funcionando correctamente"})
		})

		api.POST("/accounts", accountHandler.Create)
		api.GET("/accounts", accountHandler.GetAll)

		api.POST("/transactions", transactionHandler.Create)
		api.GET("/transactions", transactionHandler.GetAll)
		api.GET("/transactions/balance", transactionHandler.GetBalance)
		api.GET("/transactions/balance-by-account", transactionHandler.GetBalanceByAccount)
	}

	return router
}
