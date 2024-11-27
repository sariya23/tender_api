package routes

import "github.com/gin-gonic/gin"

func BindRoutes(r *gin.RouterGroup) {
	bind := r.Group("/bind")
	{
		bind.POST("/new", func(ctx *gin.Context) { ctx.JSON(200, gin.H{"message": "post at /api/bind/new"}) })
	}
}
