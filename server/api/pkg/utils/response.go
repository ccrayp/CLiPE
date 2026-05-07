package utils

import (
	"time"

	"github.com/gin-gonic/gin"
)

type Common struct {
	Success      bool   `json:"success"`
	Message      any    `json:"message"`
	Status       int    `json:"status"`
	Endpoint     string `json:"endpoint"`
	ResponseTime string `json:"response_time"`
}

type Success struct {
	Common
	Data any `json:"data"`
}

type Error struct {
	Common
	Error string `json:"error"`
}

func RespondSuccess(ctx *gin.Context, httpStatus int, message any, data any) {
	ctx.JSON(httpStatus, Success{
		Common: Common{
			Success:      true,
			Message:      message,
			Status:       httpStatus,
			Endpoint:     ctx.FullPath(),
			ResponseTime: time.Now().Local().Format("2006-01-02 15:04:05 MST"),
		},
		Data: data,
	})
}

func RespondError(ctx *gin.Context, httpStatus int, message any, error string) {
	ctx.JSON(httpStatus, Error{
		Common: Common{
			Success:      false,
			Message:      message,
			Status:       httpStatus,
			Endpoint:     ctx.FullPath(),
			ResponseTime: time.Now().Local().Format("2006-01-02 15:04:05 MST"),
		},
		Error: error,
	})
}
