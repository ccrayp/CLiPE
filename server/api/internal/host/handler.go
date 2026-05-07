package host

import (
	"clipe/internal/auth"
	"clipe/pkg/utils"
	"encoding/json"
	"net"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type HostHandler struct {
	repository_ *HostRepository
	debug_      bool
}

func NewHostHandler(repo *HostRepository, debug bool) *HostHandler {
	return &HostHandler{
		repository_: repo,
		debug_:      debug,
	}
}

func (h *HostHandler) Filter(ctx *gin.Context) {

	if auth.Require(ctx, auth.User, auth.Installer) == nil {
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

	var filter HostDTO
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
		"hosts":  data,
		"limit":  limit,
		"offset": offset,
		"count":  count,
	})
}

func (h *HostHandler) Create(ctx *gin.Context) {

	if auth.Require(ctx, auth.User, auth.Installer) == nil {
		return
	}

	var dto CreateHostDTO

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

	if net.ParseIP(dto.IP) == nil {
		utils.RespondError(ctx, http.StatusBadRequest, "invalid ip", "invalid ip")
		return
	}

	utils.RespondSuccess(ctx, http.StatusCreated, nil, gin.H{
		"id": id,
	})
}

func (h *HostHandler) Update(ctx *gin.Context) {

	if auth.Require(ctx, auth.User) == nil {
		return
	}

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "invalid id", err.Error())
		return
	}

	var dto CreateHostDTO

	decoder := json.NewDecoder(ctx.Request.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&dto); err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "invalid body", err.Error())
		return
	}

	if net.ParseIP(dto.IP) == nil {
		utils.RespondError(ctx, http.StatusBadRequest, "invalid ip", "invalid ip")
		return
	}

	err = h.repository_.Update(uint(id), &dto)
	if err != nil {
		utils.RespondError(ctx, http.StatusInternalServerError, nil, err.Error())
		return
	}

	utils.RespondSuccess(ctx, http.StatusOK, nil, nil)
}

func (h *HostHandler) Delete(ctx *gin.Context) {

	if auth.Require(ctx, auth.User, auth.Installer) == nil {
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
