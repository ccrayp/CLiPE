package aggregator

import (
	"clipe/pkg/utils"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Aggregator struct {
	repository_ *AggreagtorRepository
}

func NewAggregator(repository *AggreagtorRepository) *Aggregator {
	return &Aggregator{
		repository_: repository,
	}
}

func (a *Aggregator) Get(ctx *gin.Context) {
	var filter AggregatorDTO

	decoder := json.NewDecoder(ctx.Request.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&filter); err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "invalid body", err.Error())
		return
	}

	user_id, err := a.repository_.FindUserIdByName(filter.UserName)
	if err != nil {
		utils.RespondError(ctx, http.StatusInternalServerError, "failed to find user", err.Error())
		return
	}

	host_id, err := a.repository_.FinHostIdByIp(filter.HostIp)
	if err != nil {
		utils.RespondError(ctx, http.StatusInternalServerError, "failed to find host", err.Error())
		return
	}

	service_id, err := a.repository_.FindServiceIdByName(filter.ServiceName)
	if err != nil {
		utils.RespondError(ctx, http.StatusInternalServerError, "failed to find service", err.Error())
		return
	}

	action_id, err := a.repository_.FindActionIdByName(filter.ActionName)
	if err != nil {
		utils.RespondError(ctx, http.StatusInternalServerError, "failed to find action", err.Error())
		return
	}

	if user_id == nil || host_id == nil || service_id == nil || action_id == nil {
		utils.RespondError(ctx, http.StatusInternalServerError, "internal aggregation error: nil ids", "internal aggregation error: nil ids")
		return
	}

	policyData, err := a.repository_.FindPolicy(*user_id, *host_id, *service_id, *action_id)
	if err != nil {
		utils.RespondError(ctx, http.StatusInternalServerError, "failed to find policy", err.Error())
		return
	}

	if policyData == nil {
		utils.RespondError(ctx, http.StatusNotFound, "policy not found", "policy not found")
		return
	}

	if policyData.RuleID == nil {
		utils.RespondError(ctx, http.StatusInternalServerError, "policy has no rule", "policy has no rule")
		return
	}

	ruleData, err := a.repository_.FindRuleById(*policyData.RuleID)
	if err != nil {
		utils.RespondError(ctx, http.StatusInternalServerError, "failed to find rule", err.Error())
		return
	}

	if ruleData == nil {
		utils.RespondError(ctx, http.StatusInternalServerError, "rule is empty", "rule is empty")
		return
	}

	var conditions []Condition
	if err := json.Unmarshal(ruleData.Condition, &conditions); err != nil {
		utils.RespondError(ctx, http.StatusInternalServerError, "invalid rule condition format", err.Error())
		return
	}

	response := PolicyMatchResponse{
		Policy: PolicyResponse{
			ID:        policyData.PolicyID,
			Name:      policyData.PolicyName,
			UserID:    *policyData.UserID,
			HostID:    *policyData.HostID,
			ServiceID: *policyData.ServiceID,
			ActionID:  *policyData.ActionID,
			Status:    policyData.Status,
		},
		Rule: Rule{
			Conditions: conditions,
			Effect:     ruleData.Effect,
		},
	}

	utils.RespondSuccess(ctx, http.StatusOK, nil, response)
}
