package auth

import (
	"clipe/pkg/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

var (
	User           = SubjectRule{Type: PrincipalHuman}
	Installer      = SubjectRule{Type: PrincipalMachine, ID: "installer"}
	DecisionServer = SubjectRule{Type: PrincipalMachine, ID: "decision_server"}
)

type SubjectRule struct {
	Type PrincipalType
	ID   string
}

func Require(ctx *gin.Context, allowed ...SubjectRule) *Principal {
	principal := GetPrincipal(ctx)
	if principal == nil {
		utils.RespondError(ctx, http.StatusUnauthorized, "unauthorized", "unauthorized")
		ctx.Abort()
		return nil
	}

	for _, rule := range allowed {
		if principal.Type != rule.Type {
			continue
		}

		if rule.ID == "" {
			return principal
		}

		if principal.ID == rule.ID {
			return principal
		}
	}

	utils.RespondError(ctx, http.StatusForbidden, "forbidden", "forbidden")
	ctx.Abort()
	return nil
}
