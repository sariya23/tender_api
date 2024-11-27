package routes

import "github.com/gin-gonic/gin"

func TenderRoutes(r *gin.RouterGroup) {
	tender := r.Group("/tender")
	{
		tender.POST("/new", func(ctx *gin.Context) { ctx.JSON(200, gin.H{"message": "post at /api/tender/new"}) })
	}
}
