package middlewares

import (
	"projet_vdfi/initializers"

	"github.com/gin-gonic/gin"
)

// Middleware pour injecter la base de données dans le contexte
func DBMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("db", initializers.DB)
		c.Next()
	}
}
