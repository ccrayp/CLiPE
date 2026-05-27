package service

import (
	"clipe/internal/auth"
	"clipe/pkg/utils"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ServiceHandler struct {
	repository_ *ServiceRepository
	debug_      bool
}

func NewServiceHandler(repo *ServiceRepository, debug bool) *ServiceHandler {
	return &ServiceHandler{
		repository_: repo,
		debug_:      debug,
	}
}

func (h *ServiceHandler) Filter(ctx *gin.Context) {

	if auth.Require(ctx, auth.User, auth.DecisionServer) == nil {
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

	var filter ServiceDTO
	decoder := json.NewDecoder(ctx.Request.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&filter); err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Неверные параметры запроса", "Неверные параметры запроса")
		return
	}

	data, err := h.repository_.Select(&filter, limit, offset)
	if err != nil {
		utils.RespondError(ctx, http.StatusInternalServerError, nil, "Ошибка при запросе данных")
		return
	}

	count, err := h.repository_.Count()
	if err != nil {
		utils.RespondError(ctx, http.StatusInternalServerError, "Ошибка при подсчёте данных", "Ошибка при подсчёте данных")
		return
	}

	utils.RespondSuccess(ctx, http.StatusOK, nil, gin.H{
		"services": data,
		"limit":    limit,
		"offset":   offset,
		"count":    count,
	})
}

func (h *ServiceHandler) Create(ctx *gin.Context) {

	if auth.Require(ctx, auth.User, auth.DecisionServer) == nil {
		return
	}

	var dto CreateServiceDTO

	decoder := json.NewDecoder(ctx.Request.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&dto); err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Неверные параметры запроса", err.Error())
		return
	}

	if len(dto.ServiceName) <= 0 || len(dto.ServiceName) > 100 {
		utils.RespondError(ctx, http.StatusBadRequest, "Неверное имя сервиса (не более 100 символов)", "Неверное имя сервиса (не более 100 символов)")
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

func (h *ServiceHandler) Update(ctx *gin.Context) {

	if auth.Require(ctx, auth.User) == nil {
		return
	}

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Неверный идентификатор", "Неверный идентификатор")
		return
	}

	var dto CreateServiceDTO

	decoder := json.NewDecoder(ctx.Request.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&dto); err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Неверные параметры запроса", "Неверные параметры запроса")
		return
	}

	if len(dto.ServiceName) <= 0 || len(dto.ServiceName) > 100 {
		utils.RespondError(ctx, http.StatusBadRequest, "Неверное имя сервиса (не более 100 символов)", "Неверное имя сервиса (не более 100 символов)")
		return
	}

	err = h.repository_.Update(uint(id), &dto)
	if err != nil {
		utils.RespondError(ctx, http.StatusInternalServerError, nil, "Ошибка при обновлении записи")
		return
	}

	utils.RespondSuccess(ctx, http.StatusOK, nil, nil)
}

func (h *ServiceHandler) Delete(ctx *gin.Context) {

	if auth.Require(ctx, auth.User) == nil {
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
