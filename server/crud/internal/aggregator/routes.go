package aggregator

import (
	"clipe/pkg/database"

	"github.com/gin-gonic/gin"
)

func InitRoutes(r *gin.RouterGroup, db *database.DB, debug bool) {
	repository := NewAggregatorRepository(db)
	aggregator := NewAggregator(repository, debug)

	r.POST("/aggregator", aggregator.Get)
}
