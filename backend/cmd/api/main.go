package main

import (
	"cartera-app/backend/internal/config"
	"cartera-app/backend/internal/routes"
	"log"
)

func main() {
	cfg := config.LoadConfig()
	router := routes.SetupRouter()

	log.Printf("Servidor corriendo en http://localhost:%s", cfg.AppPort)

	if err := router.Run(":" + cfg.AppPort); err != nil {
		log.Fatal("No se pudo iniciar el servidor: ", err)
	}
}
