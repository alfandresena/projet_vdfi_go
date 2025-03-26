package middlewares

import (
	"net/http"
	"projet_vdfi/initializers"
	"projet_vdfi/models"

	"github.com/gin-gonic/gin"
)

// Middleware pour vérifier si l'utilisateur est un administrateur
func GetAuthenticatedAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Récupérer l'ID utilisateur à partir du contexte (défini par AuthMiddleware)
		userID, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		// Vérifier si userID est bien de type uint
		uid, ok := userID.(uint)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID"})
			c.Abort()
			return
		}

		// Récupérer l'utilisateur depuis la base de données
		var user models.User
		result := initializers.DB.First(&user, uid)
		if result.Error != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
			c.Abort()
			return
		}

		// Vérifier si l'utilisateur est un administrateur
		if !user.IsAdmin {
			c.JSON(http.StatusForbidden, gin.H{"error": "Accès refusé, administrateur requis"})
			c.Abort()
			return
		}

		// Si tout est bon, continuer la requête
		c.Next()
	}
}
