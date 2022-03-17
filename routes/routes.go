package routes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var DB *gorm.DB

type routes struct {
	router *gin.Engine
}

func NewRoutes(databaseConnection *gorm.DB) routes {
	r := routes{
		router: gin.New(),
	}

	DB = databaseConnection
	v1 := r.router.Group("/v1")

	r.addPing(v1)
	r.addAuthentication(v1)

	return r
}

func (r routes) Run(addr ...string) error {
	return r.router.Run()
}
