package auth

import (
	"clipe/pkg/database"
	"clipe/pkg/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	repo *Repository
}

func NewHandler(db *database.DB) *Handler {
	return &Handler{
		repo: NewRepository(db),
	}
}

func (h *Handler) Login(ctx *gin.Context) {
	var req LoginRequest

	if err := ctx.BindJSON(&req); err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error(), err.Error())
		return
	}

	user, err := h.repo.GetUserByUsername(req.Username)
	if err != nil {
		utils.RespondError(ctx, http.StatusUnauthorized, "wrong login or password", "wrong login or password")
		return
	}

	if !CheckPassword(req.Password, user.Password) {
		utils.RespondError(ctx, http.StatusUnauthorized, "wrong login or password", "wrong login or password")
		return
	}

	accessToken, err := GenerateToken(user.Username)
	if err != nil {
		utils.RespondError(ctx, http.StatusInternalServerError, err.Error(), err.Error())
		return
	}

	refreshToken, err := GenerateRefreshToken(h.repo, user.Username)
	if err != nil {
		utils.RespondError(ctx, http.StatusInternalServerError, err.Error(), err.Error())
		return
	}

	utils.RespondSuccess(ctx, http.StatusOK, "login success", gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

func (h *Handler) Refresh(ctx *gin.Context) {
	var req RefreshRequest

	if err := ctx.BindJSON(&req); err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error(), err.Error())
		return
	}

	accessToken, refreshToken, err := RefreshAccessToken(h.repo, req.RefreshToken)
	if err != nil {
		utils.RespondError(ctx, http.StatusUnauthorized, "invalid refresh token", "invalid refresh token")
		return
	}

	utils.RespondSuccess(ctx, http.StatusOK, "successfully refreshed", gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

func (h *Handler) Logout(ctx *gin.Context) {
	var req LogoutRequest

	if err := ctx.BindJSON(&req); err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error(), err.Error())
		return
	}

	if req.RefreshToken == "" {
		utils.RespondError(ctx, http.StatusBadRequest, "refresh token required", "refresh token required")
		return
	}

	token, err := h.repo.GetRefreshToken(req.RefreshToken)
	if err != nil {
		utils.RespondError(ctx, http.StatusInternalServerError, "error while token existing checking", err.Error())
		return
	}

	if token == nil {
		utils.RespondError(ctx, http.StatusInternalServerError, "token was not found", "token was not found")
		return
	}

	err = h.repo.DeleteRefreshToken(req.RefreshToken)
	if err != nil {
		utils.RespondError(ctx, http.StatusInternalServerError, err.Error(), err.Error())
		return
	}

	utils.RespondSuccess(ctx, http.StatusOK, "logout success", nil)
}

func (h *Handler) Hash(ctx *gin.Context) {
	password := ctx.Query("password")
	if password == "" {
		utils.RespondError(ctx, http.StatusBadRequest, "password query is required", "password query is required")
		return
	}

	hash, err := utils.HashPassword(password)
	if err != nil {
		utils.RespondError(ctx, http.StatusInternalServerError, err.Error(), err.Error())
		return
	}

	utils.RespondSuccess(ctx, http.StatusOK, "password hash generated", gin.H{
		"password":      password,
		"password_hash": hash,
	})
}
