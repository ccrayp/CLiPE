package policy

import (
	"clipe/internal/auth"
	"clipe/pkg/utils"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PolicyHandler struct {
	repository_ *PolicyRepository
	debug_      bool
}

func NewPolicyHandler(service *PolicyRepository, debug bool) *PolicyHandler {
	return &PolicyHandler{
		repository_: service,
		debug_:      debug,
	}
}

func (h *PolicyHandler) Filter(ctx *gin.Context) {

	if auth.Require(ctx, auth.User) == nil {
		return
	}

	temp := ctx.Query("limit")
	limit, err := strconv.Atoi(temp)
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Неверный лимит", "Неверный лимит")
		return
	}

	temp = ctx.Query("offset")
	offset, err := strconv.Atoi(temp)
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Неверный сдвиг", "Неверный сдвиг")
		return
	}

	var policy PolicyDTO
	decoder := json.NewDecoder(ctx.Request.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&policy); err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Неверные параметры запроса", "Неверные параметры запроса")
		return
	}

	if h.debug_ {
		fmt.Printf("search policy: %v", policy)
	}

	data, err := h.repository_.Select(&policy, limit, offset)
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
		"policies": data,
		"limit":    limit,
		"offset":   offset,
		"count":    count,
	})
}

func (h *PolicyHandler) Create(ctx *gin.Context) {

	if auth.Require(ctx, auth.User) == nil {
		return
	}

	var policy CreatePolicyDTO

	decoder := json.NewDecoder(ctx.Request.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&policy); err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Неверные параметры запроса", "Неверные параметры запроса")
		return
	}

	if len(policy.PolicyName) <= 0 || len(policy.PolicyName) > 100 {
		utils.RespondError(ctx, http.StatusBadRequest, "Неверная длина названия политики (не более 100 символов)", "Неверная длина названия политики (не более 100 символов)")
		return
	}

	data, err := h.repository_.Create(&policy)
	if err != nil {
		utils.RespondError(ctx, http.StatusInternalServerError, nil, "Ошибка при создании записи")
		return
	}

	utils.RespondSuccess(ctx, http.StatusCreated, nil, gin.H{
		"id": data,
	})
}

func (h *PolicyHandler) Update(ctx *gin.Context) {

	if auth.Require(ctx, auth.User) == nil {
		return
	}

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Неверный идентификатор", "Неверный идентификатор")
	}

	var policy CreatePolicyDTO
	decoder := json.NewDecoder(ctx.Request.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&policy); err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Неверные параметры запроса", "Неверные параметры запроса")
		return
	}

	if len(policy.PolicyName) <= 0 || len(policy.PolicyName) > 100 {
		utils.RespondError(ctx, http.StatusBadRequest, "Неверное название политики (не более 100 символов)", "Неверное название политики (не более 100 символов)")
		return
	}

	err = h.repository_.Update(uint(id), &policy)
	if err != nil {
		utils.RespondError(ctx, http.StatusInternalServerError, nil, "Ошибка при обновлении записи")
		return
	}

	utils.RespondSuccess(ctx, http.StatusCreated, nil, nil)
}

func (h *PolicyHandler) Delete(ctx *gin.Context) {

	if auth.Require(ctx, auth.User) == nil {
		return
	}

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Неверный идентификатор", "Неверный идентификатор")
	}

	err = h.repository_.Delete(uint(id))
	if err != nil {
		utils.RespondError(ctx, http.StatusInternalServerError, nil, "Ошибка при удалении записи")
	}

	utils.RespondSuccess(ctx, http.StatusOK, nil, nil)
}
