package routes

import (
	"tender/internal/handlers/ping"

	"github.com/gin-gonic/gin"
)

func PingRoutes(r *gin.RouterGroup) {
	bind := r.Group("/ping")
	{
		bind.GET("/", ping.Ping)
	}
}
