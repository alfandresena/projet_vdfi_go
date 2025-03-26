package initializers

import (
	"log"

	"github.com/joho/godotenv"
)

func LoadEnvVariables() {
	// Charger les variables d'environnement
	err := godotenv.Load()
	if err != nil {
		log.Println("Erreur lors du chargement du fichier .env")
	} else {
		log.Println("Fichier .env chargé avec succès")
	}
}
