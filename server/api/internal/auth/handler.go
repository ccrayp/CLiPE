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
		utils.RespondError(ctx, http.StatusUnauthorized, "Неверный логин или пароль", "Неверный логин или пароль")
		return
	}

	if !CheckPassword(req.Password, user.Password) {
		utils.RespondError(ctx, http.StatusUnauthorized, "Неверный логин или пароль", "Неверный логин или пароль")
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

	utils.RespondSuccess(ctx, http.StatusOK, "Успешный вход в систему", gin.H{
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
		utils.RespondError(ctx, http.StatusUnauthorized, "Неверный токен обновления", "Неверный токен обновления")
		return
	}

	utils.RespondSuccess(ctx, http.StatusOK, "Токены успешно обновлены", gin.H{
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
		utils.RespondError(ctx, http.StatusBadRequest, "Необходим токен обновления", "Необходим токен обновления")
		return
	}

	token, err := h.repo.GetRefreshToken(req.RefreshToken)
	if err != nil {
		utils.RespondError(ctx, http.StatusInternalServerError, "Ошибка при проверке существования токена", err.Error())
		return
	}

	if token == nil {
		utils.RespondError(ctx, http.StatusInternalServerError, "Токен не был найден", "Токен не был найден")
		return
	}

	err = h.repo.DeleteRefreshToken(req.RefreshToken)
	if err != nil {
		utils.RespondError(ctx, http.StatusInternalServerError, err.Error(), "Ошибка при удалении токена обновления")
		return
	}

	utils.RespondSuccess(ctx, http.StatusOK, "Выход успешно выполнен", "Выход успешно выполнен")
}

func (h *Handler) Hash(ctx *gin.Context) {
	password := ctx.Query("password")
	if password == "" {
		utils.RespondError(ctx, http.StatusBadRequest, "Необходим пароль в параметрах запроса", "Необходим пароль в параметрах запроса")
		return
	}

	hash, err := utils.HashPassword(password)
	if err != nil {
		utils.RespondError(ctx, http.StatusInternalServerError, err.Error(), "Ошибка при хэшировании пароля")
		return
	}

	utils.RespondSuccess(ctx, http.StatusOK, "Хэш пароля успешно создан", gin.H{
		"password":      password,
		"password_hash": hash,
	})
}
