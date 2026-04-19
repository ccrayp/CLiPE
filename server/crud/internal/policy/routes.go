package policy

import (
	"clipe/pkg/database"

	"github.com/gin-gonic/gin"
)

func InitRoutes(r *gin.RouterGroup, db *database.DB, debug bool) {
	repository := NewPolicyRep(db)
	handler := NewPolicyHandler(repository, debug)

	group := r.Group("/policies")

	group.GET("", handler.Filter)
	group.POST("", handler.Create)
	group.PUT("/:id", handler.Update)
	group.DELETE("/:id", handler.Delete)
}
