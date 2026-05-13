package policycontent

import (
	"clipe/internal/auth"
	"clipe/pkg/utils"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PolicyContentHandler struct {
	repository_ *PolicyContentRepository
	debug_      bool
}

func NewPolicyContentHandler(repository *PolicyContentRepository, debug bool) *PolicyContentHandler {
	return &PolicyContentHandler{
		repository_: repository,
		debug_:      debug,
	}
}

func (h *PolicyContentHandler) Filter(ctx *gin.Context) {
	if auth.Require(ctx, auth.User) == nil {
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

	var filter PolicyContentDTO

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
		utils.RespondError(ctx, http.StatusInternalServerError, "Ошибка при подсчёте результатов", "Ошибка при подсчёте результатов")
		return
	}

	utils.RespondSuccess(ctx, http.StatusOK, nil, gin.H{
		"policy_contents": data,
		"limit":           limit,
		"offset":          offset,
		"count":           count,
	})
}

func (h *PolicyContentHandler) Create(ctx *gin.Context) {
	if auth.Require(ctx, auth.User) == nil {
		return
	}

	var dto CreatePolicyContentDTO

	decoder := json.NewDecoder(ctx.Request.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&dto); err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Неверные параметры запроса", "Неверные параметры запроса")
		return
	}

	if err := h.repository_.Create(&dto); err != nil {
		utils.RespondError(ctx, http.StatusInternalServerError, nil, "Ошибка при создании записи")
		return
	}

	utils.RespondSuccess(ctx, http.StatusCreated, nil, nil)
}

func (h *PolicyContentHandler) Update(ctx *gin.Context) {
	if auth.Require(ctx, auth.User) == nil {
		return
	}

	policyID, err := strconv.Atoi(ctx.Param("policy_id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Неверный идентификатор политики", "Неверный идентификатор политики")
		return
	}

	serviceID, err := strconv.Atoi(ctx.Param("service_id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Неверный идентификатор сервиса", "Неверный идентификатор сервиса")
		return
	}

	var dto CreatePolicyContentDTO

	decoder := json.NewDecoder(ctx.Request.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&dto); err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Неверные параметры запроса", "Неверные параметры запроса")
		return
	}

	dto.PolicyID = uint(policyID)
	dto.ServiceID = uint(serviceID)

	if err := h.repository_.Update(uint(policyID), uint(serviceID), &dto); err != nil {
		utils.RespondError(ctx, http.StatusInternalServerError, nil, "Ошибка при обновлении записи")
		return
	}

	utils.RespondSuccess(ctx, http.StatusOK, nil, nil)
}

func (h *PolicyContentHandler) Delete(ctx *gin.Context) {
	if auth.Require(ctx, auth.User) == nil {
		return
	}

	policyID, err := strconv.Atoi(ctx.Param("policy_id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Неверный идентификатор политики", "Неверный идентификатор политики")
		return
	}

	serviceID, err := strconv.Atoi(ctx.Param("service_id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Неверный идентификатор сервиса", "Неверный идентификатор сервиса")
		return
	}

	if err := h.repository_.Delete(uint(policyID), uint(serviceID)); err != nil {
		utils.RespondError(ctx, http.StatusInternalServerError, nil, "Ошибка при удалении записи")
		return
	}

	utils.RespondSuccess(ctx, http.StatusOK, nil, nil)
}
