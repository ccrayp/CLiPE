// Выбор всех записей с фильтрами
// Выбор записи по id

package policy

import (
	"clipe/pkg/utils"
	"encoding/json"
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

	temp := ctx.Query("limit")
	limit, err := strconv.Atoi(temp)
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "invalid limit", err.Error())
		return
	}

	temp = ctx.Query("offset")
	offset, err := strconv.Atoi(temp)
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "invalid offset", err.Error())
		return
	}

	var policy PolicyDTO
	decoder := json.NewDecoder(ctx.Request.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&policy); err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "invalid body", err.Error())
		return
	}

	data, err := h.repository_.Select(&policy, limit, offset)
	if err != nil {
		utils.RespondError(ctx, http.StatusInternalServerError, nil, err.Error())
		return
	}

	utils.RespondSuccess(ctx, http.StatusOK, nil, gin.H{
		"policies": data,
	})
}

func (h *PolicyHandler) Create(ctx *gin.Context) {

	var policy CreatePolicyDTO

	decoder := json.NewDecoder(ctx.Request.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&policy); err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "invalid body", err.Error())
		return
	}

	data, err := h.repository_.Create(&policy)
	if err != nil {
		utils.RespondError(ctx, http.StatusInternalServerError, nil, err.Error())
		return
	}

	utils.RespondSuccess(ctx, http.StatusCreated, nil, gin.H{
		"id": data,
	})
}

func (h *PolicyHandler) Update(ctx *gin.Context) {

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "invalid id", err.Error())
	}

	var policy CreatePolicyDTO
	decoder := json.NewDecoder(ctx.Request.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&policy); err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "invalid body", err.Error())
		return
	}

	err = h.repository_.Update(uint(id), &policy)
	if err != nil {
		utils.RespondError(ctx, http.StatusInternalServerError, nil, err.Error())
		return
	}

	utils.RespondSuccess(ctx, http.StatusCreated, nil, nil)
}

func (h *PolicyHandler) Delete(ctx *gin.Context) {

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "invalid id", err.Error())
	}

	err = h.repository_.Delete(uint(id))
	if err != nil {
		utils.RespondError(ctx, http.StatusInternalServerError, nil, err.Error())
	}

	utils.RespondSuccess(ctx, http.StatusOK, nil, nil)
}
