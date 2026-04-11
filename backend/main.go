package main

import (
	"clipe/internal/policy"
	"clipe/pkg/config"
	"clipe/pkg/database"
	"clipe/pkg/utils"
	"fmt"
	"net/http"
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
	api := server.Group("/api")

	InitRoutes(api, db)

	server.NoRoute(func(ctx *gin.Context) {
		ctx.Redirect(http.StatusPermanentRedirect, "/api")
	})

	server.Run(fmt.Sprintf(":%s", cfg.Server.Port))
}

func InitRoutes(r *gin.RouterGroup, db *database.DB) {

	r.GET("", func(ctx *gin.Context) {
		now := time.Now().Format("2006-01-02 15:04:05")

		html := fmt.Sprintf(`
			<p>Время на сервере: %s</p>
			<p>API централизованного сервера управления доступом</p>
			<a href="/api/health">Healthcheck</a>
		`, now)

		ctx.Data(200, "text/html; charset=utf-8", []byte(html))
	})

	r.GET("/health", func(ctx *gin.Context) {
		utils.RespondSuccess(ctx, 200, nil, gin.H{
			"health": true,
		})
	})

	policy.InitRoutes(r, db)
}
