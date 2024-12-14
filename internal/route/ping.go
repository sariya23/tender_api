package route

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func PingRoute(r *gin.RouterGroup) {
	ping := r.Group("/ping")
	{
		ping.GET("/", func(c *gin.Context) {
			{
				c.JSON(http.StatusOK, struct{ Message string }{Message: "ok"})
			}
		})
	}
}
