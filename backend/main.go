package main

import (
	"clipe/internal/user"
	"clipe/pkg/config"
	"clipe/pkg/database"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	fmt.Print("Loadind .env variables...\n")

	godotenv.Load()
	cfg := config.NewConfig()

	fmt.Printf("host = %s, port = %s, name = %s, server_port = %s", cfg.Database.Host, cfg.Database.Port, cfg.Database.Name, cfg.Server.Port)
	fmt.Print(".env variables were successfully loaded!\n\n")

	fmt.Print("Initializing database connection...\n")

	db, err := database.NewDB(*cfg)

	if err != nil {
		fmt.Printf("Database connection was not successfully established...\nerror: %s", err.Error())
		return
	}

	fmt.Print("Database connection was successfully established!\n\n")

	fmt.Print("Initializing API-server...\n")

	gin.SetMode(gin.ReleaseMode)
	server := gin.Default()

	server.GET("/api", func(ctx *gin.Context) {
		now := time.Now().Format("2006-01-02 15:04:05")

		html := fmt.Sprintf(`
			<p>Время на сервере: %s</p>
			<p>API централизованного сервера управления доступом</p>
			<a href="/api/health">Healthcheck</a>
		`, now)

		ctx.Data(200, "text/html; charset=utf-8", []byte(html))
	})

	server.GET("/api/health", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"health":    "ok",
			"timestamp": time.Now().Format("2006-01-02 15:04:05"),
		})
	})

	server.GET("/api/users", func(ctx *gin.Context) {
		var users []user.User

		if err := db.Conn().Find(&users).Error; err != nil {
			ctx.JSON(500, gin.H{
				"error": err.Error(),
			})
		}

		var result []user.UserDTO
		for _, u := range users {
			result = append(result, user.ToDTO(u))
		}

		ctx.JSON(http.StatusOK, gin.H{
			"users": result,
		})
	})

	server.NoRoute(func(ctx *gin.Context) {
		ctx.Redirect(http.StatusPermanentRedirect, "/api")
	})

	fmt.Print("API-server was successfully started!\n\n")
	server.Run(fmt.Sprintf(":%s", cfg.Server.Port))
}
