package utils

import (
	"net/http"

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

type Journal struct {
	PolicyId   uint   `json:"policy_id"`
	PolicyName string `json:"policy_name"`
	RequestId  uint   `json:"request_id"`
	DecisionId uint   `json:"decision_id"`
}

type ApiResponse struct {
	Result  string  `json:"result"`
	Journal Journal `json:"journal"`
}

func GetJournal(policyId uint, policyName string, requestId, decisionId uint) Journal {
	return Journal{
		PolicyId:   policyId,
		PolicyName: policyName,
		RequestId:  requestId,
		DecisionId: decisionId,
	}
}

func AllowRespond(ctx *gin.Context, policyId uint, policyName string, requestId, decisionId uint) {
	ctx.JSON(http.StatusOK, ApiResponse{
		Result:  "ALLOW",
		Journal: GetJournal(policyId, policyName, requestId, decisionId),
	})
}

func DenyRespond(ctx *gin.Context, policyId uint, policyName string, requestId, decisionId uint) {
	ctx.JSON(http.StatusOK, ApiResponse{
		Result:  "DENY",
		Journal: GetJournal(policyId, policyName, requestId, decisionId),
	})
}

func ErrorRespond(ctx *gin.Context, httpStatus int, message any, err string) {
	ctx.JSON(httpStatus, gin.H{
		"message": message,
		"error":   err,
	})
}
