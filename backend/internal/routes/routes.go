package routes

import (
	"cartera-app/backend/internal/accounts"
	"cartera-app/backend/internal/auth"
	"cartera-app/backend/internal/config"
	"cartera-app/backend/internal/transactions"
	"cartera-app/backend/internal/transfers"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func SetupRouter(db *pgxpool.Pool, cfg *config.Config) *gin.Engine {
	router := gin.Default()

	api := router.Group("/api")
	{
		authRepo := auth.NewPostgresRepository(db)
		tokenManager := auth.NewTokenManager(cfg.JWTSecret, "cartera-app", 24*time.Hour)
		authService := auth.NewService(authRepo, tokenManager)
		authHandler := auth.NewHandler(authService)

		accountRepo := accounts.NewPostgresRepository(db)
		transactionRepo := transactions.NewPostgresRepository(db)
		transferRepo := transfers.NewPostgresRepository(db)

		accountService := accounts.NewService(accountRepo)
		transactionService := transactions.NewService(transactionRepo, accountService, accountService)

		accountHandler := accounts.NewHandler(accountService, transactionService)
		transactionHandler := transactions.NewHandler(transactionService)

		transferService := transfers.NewService(accountService, transferRepo)
		transferHandler := transfers.NewHandler(transferService)

		api.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "API funcionando correctamente"})
		})

		api.POST("/auth/register", authHandler.Register)
		api.POST("/auth/login", authHandler.Login)

		private := api.Group("")
		private.Use(auth.RequireAuth(tokenManager))
		{
			private.GET("/auth/me", authHandler.Me)

			private.POST("/accounts", accountHandler.Create)
			private.GET("/accounts", accountHandler.GetAll)
			private.GET("/accounts/:id", accountHandler.GetByID)
			private.PUT("/accounts/:id", accountHandler.Update)
			private.DELETE("/accounts/:id", accountHandler.Delete)

			private.POST("/transactions", transactionHandler.Create)
			private.GET("/transactions", transactionHandler.GetAll)
			private.GET("/transactions/balance", transactionHandler.GetBalance)
			private.GET("/transactions/balance-by-account", transactionHandler.GetBalanceByAccount)
			private.GET("/transactions/summary", transactionHandler.GetSummary)
			private.GET("/transactions/:id", transactionHandler.GetByID)
			private.PUT("/transactions/:id", transactionHandler.Update)
			private.DELETE("/transactions/:id", transactionHandler.Delete)

			private.POST("/transfers", transferHandler.Create)
		}
	}

	return router
}
