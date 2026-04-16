package service

import (
	"clipe/pkg/database"

	"github.com/gin-gonic/gin"
)

func InitRoutes(r *gin.RouterGroup, db *database.DB) {
	repository := NewServiceRep(db)
	handler := NewServiceHandler(repository)

	group := r.Group("/services")

	group.GET("", handler.Filter)
	group.POST("", handler.Create)
	group.PUT("/:id", handler.Update)
	group.DELETE("/:id", handler.Delete)
}
