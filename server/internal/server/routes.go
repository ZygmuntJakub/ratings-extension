package server

import (
	"github.com/ZygmuntJakub/mkino-extension/internal/ratings"
	"github.com/gin-gonic/gin"
)

func registerRoutes(r *gin.Engine) {
	v1 := r.Group("/api/v1")
	v1.GET("health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "OK",
		})
	})

	ratings.RatingsRoutes(v1)
}
