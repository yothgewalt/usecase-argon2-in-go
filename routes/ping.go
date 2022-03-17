package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (r routes) addPing(rg *gin.RouterGroup) {
	ping := rg.Group("/ping")
	ping.GET("/", recieverPing)
}

func recieverPing(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "You can reach the first destination for testing.",
	})
}
