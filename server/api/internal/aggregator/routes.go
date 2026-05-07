package aggregator

import (
	"clipe/internal/auth"
	"clipe/pkg/database"

	"github.com/gin-gonic/gin"
)

func InitRoutes(r *gin.RouterGroup, db *database.DB, debug bool) {
	repository := NewAggregatorRepository(db)
	aggregator := NewAggregator(repository, debug)

	r.POST("/aggregator", auth.AuthMiddleware(), aggregator.Get)
}
