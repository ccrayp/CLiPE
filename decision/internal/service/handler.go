package service

import (
	"decision/internal/model"
	"decision/pkg/utils"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	decider *Decider
}

func NewHandler(decider *Decider) *Handler {
	return &Handler{
		decider: decider,
	}
}

func (h *Handler) Decide(ctx *gin.Context) {
	var dto model.ApiRequest

	decoder := json.NewDecoder(ctx.Request.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&dto); err != nil {
		utils.ErrorRespond(ctx, http.StatusBadRequest, "invalid body", err.Error())
		return
	}

	decision, err := h.decider.Evaluate(&dto)
	if err != nil {
		utils.ErrorRespond(ctx, http.StatusInternalServerError, nil, err.Error())
		return
	}

	if decision.Result {
		utils.AllowRespond(ctx, decision.Policy.Id, decision.Policy.Name, decision.RequestId, decision.DecisionId)
		return
	} else {
		utils.DenyRespond(ctx, decision.Policy.Id, decision.Policy.Name, decision.RequestId, decision.DecisionId)
		return
	}
}
