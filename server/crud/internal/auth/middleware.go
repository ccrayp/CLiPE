package auth

import (
	"clipe/pkg/utils"
	"errors"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		internalToken := ctx.GetHeader("X-Internal-Token")
		caller := ctx.GetHeader("X-Caller")

		switch {
		case authHeader != "" && internalToken != "":
			utils.RespondError(ctx, http.StatusUnauthorized, "ambiguous auth method", "ambiguous auth method")
			ctx.Abort()
			return

		case authHeader != "":
			claims, err := ValidateToken(authHeader)
			if err != nil {
				utils.RespondError(ctx, http.StatusUnauthorized, err.Error(), err.Error())
				ctx.Abort()
				return
			}

			ctx.Set("principal", BuildPrincipalFromClaims(claims))

		case internalToken != "":
			principal, err := TokenAuth(caller, internalToken)
			if err != nil {
				utils.RespondError(ctx, http.StatusUnauthorized, err.Error(), err.Error())
				ctx.Abort()
				return
			}

			ctx.Set("principal", principal)

		default:
			utils.RespondError(ctx, http.StatusUnauthorized, "need JWT or Internal token", "need JWT or Internal token")
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}

func TokenAuth(caller, token string) (*Principal, error) {
	switch caller {
	case os.Getenv("DECISION_ID"):
		if token != os.Getenv("DECISION_TOKEN") {
			return nil, errors.New("invalid internal token")
		}
		return &Principal{Type: PrincipalMachine, ID: "decision_server"}, nil

	case os.Getenv("INSTALLER_ID"):
		if token != os.Getenv("INSTALLER_TOKEN") {
			return nil, errors.New("invalid internal token")
		}
		return &Principal{Type: PrincipalMachine, ID: "installer"}, nil
	}

	return nil, errors.New("unknown caller")
}

func GetPrincipal(ctx *gin.Context) *Principal {
	val, ok := ctx.Get("principal")
	if !ok {
		return nil
	}

	p, ok := val.(*Principal)
	if !ok {
		return nil
	}

	return p
}
