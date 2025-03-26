package main

import (
	"fmt"
	controllers "projet_vdfi/controller"
	"projet_vdfi/initializers"
	"projet_vdfi/middlewares"

	"github.com/gin-gonic/gin"
)

func init() {
	// Initialisation des variables d'environnement
	initializers.LoadEnvVariables()

	// Connexion à la base de données
	initializers.ConnectToDb()

	// Synchronisation de la base de données
	initializers.SyncDatabase()

	fmt.Println("Initialisation terminée...")
}

func main() {

	// Démarrer le serveur Gin
	r := gin.Default()

	/***********************************************/

	// Route pour l'inscription
	r.POST("/signup", controllers.Signup)

	// Route pour la connexion
	r.POST("/login", controllers.Login)

	// Route pour la déconnexion
	r.GET("/logout", controllers.Logout)

	// Route pour la récupération des utilisateurs
	r.GET("/user", middlewares.AuthMiddleware(), controllers.GetUser)

	// Route pour la mise à jour des infos utilisateurs
	r.PUT("/user/update", middlewares.AuthMiddleware(), controllers.UpdateUser)

	// Route pour la suppression d'un utilisateur
	r.DELETE("/user/delete", middlewares.AuthMiddleware(), controllers.DeleteUser)

	// Route pour changer l'utilisateur en admin
	r.POST("/user/admin", middlewares.AuthMiddleware(), controllers.PromoteToAdmin)

	/***********************************************/

	// Routes des événements
	eventRoutes := r.Group("/events")
	{
		eventRoutes.GET("/", middlewares.AuthMiddleware(), controllers.GetEvents)
		eventRoutes.GET("/:id", middlewares.AuthMiddleware(), controllers.GetEventByID)
		eventRoutes.POST("/", middlewares.AuthMiddleware(), controllers.CreateEvent)      // ADMIN
		eventRoutes.PUT("/:id", middlewares.AuthMiddleware(), controllers.UpdateEvent)    // ADMIN
		eventRoutes.DELETE("/:id", middlewares.AuthMiddleware(), controllers.DeleteEvent) // ADMIN
	}

	/***********************************************/
	// Routes publiques pour récupérer les paroles
	publicLyricsRoutes := r.Group("/lyrics")
	{
		publicLyricsRoutes.GET("/", controllers.GetLyrics)        // Liste toutes les paroles
		publicLyricsRoutes.GET("/:id", controllers.GetLyricsByID) // Récupère une chanson par ID
	}

	// Routes admin pour les paroles
	adminLyricsRoutes := r.Group("/lyrics")
	adminLyricsRoutes.Use(middlewares.AuthMiddleware())        // Vérification utilisateur connecté
	adminLyricsRoutes.Use(middlewares.GetAuthenticatedAdmin()) // Vérification admin
	{
		adminLyricsRoutes.POST("/", controllers.CreateLyrics)
		adminLyricsRoutes.PUT("/:id", controllers.UpdateLyrics)
		adminLyricsRoutes.DELETE("/:id", controllers.DeleteLyrics)
	}

	/***********************************************/

	// Démarrer le serveur
	if err := r.Run(":3000"); err != nil {
		fmt.Println("❌ Erreur au démarrage du serveur : ", err)
	}
}
