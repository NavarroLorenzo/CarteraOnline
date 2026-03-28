package main

import (
	"log"

	"cartera-app/backend/internal/config"
	"cartera-app/backend/internal/database"
	"cartera-app/backend/internal/routes"
)

func main() {
	cfg := config.LoadConfig()
	db := database.NewPool(cfg)
	defer db.Close()

	router := routes.SetupRouter(db, cfg)

	log.Printf("Servidor corriendo en http://localhost:%s", cfg.AppPort)

	if err := router.Run(":" + cfg.AppPort); err != nil {
		log.Fatal("No se pudo iniciar el servidor: ", err)
	}
}
