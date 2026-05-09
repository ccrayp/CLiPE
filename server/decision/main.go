package main

import (
	"decision/internal/client"
	"decision/internal/service"
	"decision/pkg/config"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {

	godotenv.Load()
	cfg := config.NewConfig()

	server := gin.Default()
	api := server.Group("")

	client := client.NewClient(cfg)

	InitRoutes(api, client, cfg.Server.ApiVersion, cfg.DefaultDecision)

	server.Run(":8080")
}

func InitRoutes(r *gin.RouterGroup, client *client.Client, apiVersion string, defaultDecision bool) {

	r.GET("/", func(ctx *gin.Context) {
		now := time.Now().Format("2006-01-02 15:04:05")

		html := fmt.Sprintf(`
			<p>Время на сервере: %s</p>
			<p>API централизованного сервера управления доступом</p>
			<a href="/api/v%s/health">Healthcheck</a>
		`, now, apiVersion)

		ctx.Data(200, "text/html; charset=utf-8", []byte(html))
	})

	r.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"health": true,
			"apiUrl": client.CheckApiUrl(),
		})
	})

	decider := service.NewDecider(client, defaultDecision)
	handler := service.NewHandler(decider)

	r.POST("/decide", handler.Decide)
}
