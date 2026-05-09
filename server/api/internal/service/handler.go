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
		utils.RespondError(ctx, http.StatusBadRequest, "invalid limit", err.Error())
		return
	}

	offset, err := strconv.Atoi(ctx.Query("offset"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "invalid offset", err.Error())
		return
	}

	var filter ServiceDTO
	decoder := json.NewDecoder(ctx.Request.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&filter); err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "invalid body", err.Error())
		return
	}

	data, err := h.repository_.Select(&filter, limit, offset)
	if err != nil {
		utils.RespondError(ctx, http.StatusInternalServerError, nil, err.Error())
		return
	}

	count, err := h.repository_.Count()
	if err != nil {
		utils.RespondError(ctx, http.StatusInternalServerError, "error while get count", err.Error())
		return
	}

	utils.RespondSuccess(ctx, http.StatusOK, nil, gin.H{
		"users":  data,
		"limit":  limit,
		"offset": offset,
		"count":  count,
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
		utils.RespondError(ctx, http.StatusBadRequest, "invalid body", err.Error())
		return
	}

	id, err := h.repository_.Create(&dto)
	if err != nil {
		utils.RespondError(ctx, http.StatusInternalServerError, nil, err.Error())
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
		utils.RespondError(ctx, http.StatusBadRequest, "invalid id", err.Error())
		return
	}

	var dto CreateServiceDTO

	decoder := json.NewDecoder(ctx.Request.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&dto); err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "invalid body", err.Error())
		return
	}

	err = h.repository_.Update(uint(id), &dto)
	if err != nil {
		utils.RespondError(ctx, http.StatusInternalServerError, nil, err.Error())
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
		utils.RespondError(ctx, http.StatusBadRequest, "invalid id", err.Error())
		return
	}

	err = h.repository_.Delete(uint(id))
	if err != nil {
		utils.RespondError(ctx, http.StatusInternalServerError, nil, err.Error())
		return
	}

	utils.RespondSuccess(ctx, http.StatusOK, nil, nil)
}
