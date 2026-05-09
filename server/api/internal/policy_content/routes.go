package policycontent

import (
	"clipe/internal/auth"
	"clipe/pkg/database"

	"github.com/gin-gonic/gin"
)

func InitRoutes(r *gin.RouterGroup, db *database.DB, debug bool) {
	repository := NewPolicyContentRep(db)
	handler := NewPolicyContentHandler(repository, debug)

	group := r.Group("/policy-contents")
	group.Use(auth.AuthMiddleware())

	group.POST("/search", handler.Filter)
	group.POST("", handler.Create)
	group.PUT("/:policy_id/:service_id", handler.Update)
	group.DELETE("/:policy_id/:service_id", handler.Delete)
}
