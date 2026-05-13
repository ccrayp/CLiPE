package user

import (
	"clipe/internal/auth"
	"clipe/pkg/utils"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	repository_ *UserRepository
	debug_      bool
}

func NewUserHandler(repo *UserRepository, debug bool) *UserHandler {
	return &UserHandler{
		repository_: repo,
		debug_:      debug,
	}
}

func (h *UserHandler) Filter(ctx *gin.Context) {

	if auth.Require(ctx, auth.User, auth.Installer) == nil {
		return
	}

	limit, err := strconv.Atoi(ctx.Query("limit"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Неверный лимит", "Неверный лимит")
		return
	}

	offset, err := strconv.Atoi(ctx.Query("offset"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Неверный сдвиг", "Неверный сдвиг")
		return
	}

	var filter UserDTO
	decoder := json.NewDecoder(ctx.Request.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&filter); err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Неверные параметры запроса", "Неверные параметры запроса")
		return
	}

	if h.debug_ {
		fmt.Printf("search user: %v\n", filter)
	}

	data, err := h.repository_.Select(&filter, limit, offset)
	if err != nil {
		utils.RespondError(ctx, http.StatusInternalServerError, nil, "Ошибка при запросе данных")
		return
	}

	count, err := h.repository_.Count()
	if err != nil {
		utils.RespondError(ctx, http.StatusInternalServerError, "Ошибка при подсчёте результатов", "Ошибка при подсчёте результатов")
		return
	}

	utils.RespondSuccess(ctx, http.StatusOK, nil, gin.H{
		"users":  data,
		"limit":  limit,
		"offset": offset,
		"count":  count,
	})
}

func (h *UserHandler) Create(ctx *gin.Context) {

	if auth.Require(ctx, auth.User, auth.Installer) == nil {
		return
	}

	var dto CreateUserDTO

	decoder := json.NewDecoder(ctx.Request.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&dto); err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Неверные параметры запроса", "Неверные параметры запроса")
		return
	}

	if dto.GID <= 0 || dto.UID <= 0 {
		utils.RespondError(ctx, http.StatusBadRequest, "Неверный UID или GID (должны быть больше 0)", "Неверный UID или GID (должны быть больше 0)")
		return
	}

	if len(dto.UserName) == 0 || len(dto.UserName) >= 100 {
		utils.RespondError(ctx, http.StatusBadRequest, "Неверное имя пользователя (не более 100 символов)", "Неверное имя пользователя (не более 100 символов)")
		return
	}

	id, err := h.repository_.Create(&dto)
	if err != nil {
		utils.RespondError(ctx, http.StatusInternalServerError, nil, "Ошибка при создании записи")
		return
	}

	utils.RespondSuccess(ctx, http.StatusCreated, nil, gin.H{
		"id": id,
	})
}

func (h *UserHandler) Update(ctx *gin.Context) {

	if auth.Require(ctx, auth.User) == nil {
		return
	}

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Неверный идентификатор", "Неверный идентификатор")
		return
	}

	var dto CreateUserDTO

	decoder := json.NewDecoder(ctx.Request.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&dto); err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Неверные параметры запроса", "Неверные параметры запроса")
		return
	}

	if dto.GID <= 0 || dto.UID <= 0 {
		utils.RespondError(ctx, http.StatusBadRequest, "Неверный UID или GID (должны быть больше 0)", "Неверный UID или GID (должны быть больше 0)")
		return
	}

	if len(dto.UserName) == 0 || len(dto.UserName) >= 100 {
		utils.RespondError(ctx, http.StatusBadRequest, "Неверное имя пользователя (не более 100 символов)", "Неверное имя пользователя (не более 100 символов)")
		return
	}

	err = h.repository_.Update(uint(id), &dto)
	if err != nil {
		utils.RespondError(ctx, http.StatusInternalServerError, nil, "Ошибка при обновлении записи")
		return
	}

	utils.RespondSuccess(ctx, http.StatusOK, nil, nil)
}

func (h *UserHandler) Delete(ctx *gin.Context) {

	if auth.Require(ctx, auth.User, auth.Installer) == nil {
		return
	}

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Неверный идентификатор", "Неверный идентификатор")
		return
	}

	err = h.repository_.Delete(uint(id))
	if err != nil {
		utils.RespondError(ctx, http.StatusInternalServerError, nil, "Ошибка при удалении записи")
		return
	}

	utils.RespondSuccess(ctx, http.StatusOK, nil, nil)
}
