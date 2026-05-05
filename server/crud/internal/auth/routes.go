package auth

import (
	"clipe/pkg/database"

	"github.com/gin-gonic/gin"
)

func InitRoutes(r *gin.RouterGroup, db *database.DB) {
	handler := NewHandler(db)

	group := r.Group("/auth")
	group.POST("/login", handler.Login)
	group.POST("/refresh", handler.Refresh)
	group.GET("/hash", handler.Hash)
}
