package server

import "github.com/gin-gonic/gin"

func NewServer() (r *gin.Engine) {
	r = gin.Default()

	r.Use(CORSMiddleware)

	registerRoutes(r)

	return
}
