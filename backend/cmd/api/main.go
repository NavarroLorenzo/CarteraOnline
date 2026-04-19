package main

import (
	"log"

	"cenz/backend/internal/config"
	"cenz/backend/internal/database"
	"cenz/backend/internal/routes"
)

func main() {
	cfg := config.LoadConfig()
	db := database.NewPool(cfg)
	defer db.Close()

	router := routes.SetupRouter(db, cfg)

	log.Printf("Servidor levantado en el puerto %s", cfg.AppPort)

	if err := router.Run(":" + cfg.AppPort); err != nil {
		log.Fatal("No se pudo iniciar el servidor: ", err)
	}
}
