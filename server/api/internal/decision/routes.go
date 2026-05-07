package decision

import (
	"clipe/internal/auth"
	"clipe/pkg/database"

	"github.com/gin-gonic/gin"
)

func InitRoutes(r *gin.RouterGroup, db *database.DB, debug bool) {
	repository := NewDecisionRep(db)
	handler := NewDecisionHandler(repository, debug)

	group := r.Group("/decisions")
	group.Use(auth.AuthMiddleware())

	group.POST("/search", handler.Filter)
	group.POST("", handler.Create)
	group.PUT("/:id", handler.Update)
	group.DELETE("/:id", handler.Delete)
}
