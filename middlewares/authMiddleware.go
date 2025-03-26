package middlewares

import (
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// Middleware d'authentification
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Clé secrète pour signer le token
		var SECRET_KEY = []byte(os.Getenv("JWT_SECRET"))

		// Récupérer le token depuis le cookie
		tokenString, err := c.Cookie("Authorization")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token requis"})
			c.Abort()
			return
		}

		// Vérifier et parser le token
		claims := &jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return SECRET_KEY, nil
		})

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token invalide: " + err.Error()})
			c.Abort()
			return
		}

		if !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token invalide ou expiré"})
			c.Abort()
			return
		}

		// Vérifier si le token est expiré
		if exp, ok := (*claims)["exp"].(float64); ok {
			if int64(exp) < time.Now().Unix() {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Token expiré"})
				c.Abort()
				return
			}
		}

		// Extraire l'ID utilisateur
		userID, ok := (*claims)["id"].(float64)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token invalide"})
			c.Abort()
			return
		}

		// Ajouter les infos dans le contexte Gin
		c.Set("userID", uint(userID))

		// Continuer la requête
		c.Next()
	}
}
