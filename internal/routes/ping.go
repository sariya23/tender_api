package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/sariya23/tender/internal/handlers/ping"
)

func PingRoutes(r *gin.RouterGroup) {
	bind := r.Group("/ping")
	{
		bind.GET("/", ping.Ping)
	}
}
