package decision

import (
	"clipe/pkg/database"

	"github.com/gin-gonic/gin"
)

func InitRoutes(r *gin.RouterGroup, db *database.DB) {
	repository := NewDecisionRep(db)
	handler := NewDecisionHandler(repository)

	group := r.Group("/decisions")

	group.GET("", handler.Filter)
	group.POST("", handler.Create)
	group.PUT("/:id", handler.Update)
	group.DELETE("/:id", handler.Delete)
}
