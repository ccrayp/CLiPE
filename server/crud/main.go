package main

import (
	"clipe/internal/aggregator"
	"clipe/internal/decision"
	"clipe/internal/host"
	"clipe/internal/policy"
	"clipe/internal/request"
	"clipe/internal/rule"
	"clipe/internal/user"
	"clipe/pkg/config"
	"clipe/pkg/database"
	"clipe/pkg/utils"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {

	godotenv.Load()
	cfg := config.NewConfig()

	db, err := database.NewDB(*cfg)
	if err != nil {
		return
	}

	server := gin.Default()
	api := server.Group("/api/v" + cfg.Server.ApiVersion + "/internal")

	InitRoutes(api, db, cfg.Server.ApiVersion, cfg.DebugMode)

	server.Run(":8080")
}

func InitRoutes(r *gin.RouterGroup, db *database.DB, apiVersion string, debug bool) {

	r.GET("", func(ctx *gin.Context) {
		now := time.Now().Format("2006-01-02 15:04:05")

		html := fmt.Sprintf(`
			<p>Время на сервере: %s</p>
			<p>CRUD API для БД</p>
			<a href="/api/v%s/internal/health">Healthcheck</a>
		`, now, apiVersion)

		ctx.Data(200, "text/html; charset=utf-8", []byte(html))
	})

	r.GET("/health", func(ctx *gin.Context) {
		utils.RespondSuccess(ctx, 200, nil, gin.H{
			"health": true,
		})
	})

	aggregator.InitRoutes(r, db, debug)

	decision.InitRoutes(r, db, debug)
	host.InitRoutes(r, db, debug)
	policy.InitRoutes(r, db, debug)
	request.InitRoutes(r, db, debug)
	rule.InitRoutes(r, db, debug)
	user.InitRoutes(r, db, debug)
}
