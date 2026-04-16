package aggregator

import (
	"clipe/pkg/database"

	"github.com/gin-gonic/gin"
)

func InitRoutes(r *gin.RouterGroup, db *database.DB) {
	repository := NewAggregatorRepository(db)
	aggregator := NewAggregator(repository)

	r.POST("/aggregator", aggregator.Get)
}
