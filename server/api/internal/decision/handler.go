package decision

import (
	"clipe/internal/auth"
	"clipe/pkg/utils"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type DecisionHandler struct {
	repository_ *DecisionRepository
	debug_      bool
}

func NewDecisionHandler(repo *DecisionRepository, debug bool) *DecisionHandler {
	return &DecisionHandler{
		repository_: repo,
		debug_:      debug,
	}
}

func (h *DecisionHandler) Filter(ctx *gin.Context) {

	if auth.Require(ctx, auth.User) == nil {
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

	var filter SearchDecisionDTO

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
		"decisions": data,
		"limit":     limit,
		"offset":    offset,
		"count":     count,
	})
}

func (h *DecisionHandler) Create(ctx *gin.Context) {

	if auth.Require(ctx, auth.DecisionServer) == nil {
		return
	}

	var dto CreateDecisionDTO

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

	if h.debug_ {
		fmt.Printf("new decision: %d\n", *id)
	}

	utils.RespondSuccess(ctx, http.StatusCreated, nil, gin.H{
		"id": id,
	})
}

func (h *DecisionHandler) Update(ctx *gin.Context) {

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "invalid id", err.Error())
		return
	}

	var dto CreateDecisionDTO

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

func (h *DecisionHandler) Delete(ctx *gin.Context) {

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
