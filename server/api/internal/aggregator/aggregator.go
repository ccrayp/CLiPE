package aggregator

import (
	"clipe/internal/auth"
	"clipe/pkg/utils"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Aggregator struct {
	repository_ *AggreagtorRepository
	debug_      bool
}

func NewAggregator(repository *AggreagtorRepository, debug bool) *Aggregator {
	return &Aggregator{
		repository_: repository,
		debug_:      debug,
	}
}

func (a *Aggregator) Get(ctx *gin.Context) {

	if auth.Require(ctx, auth.DecisionServer) == nil {
		return
	}

	var filter AggregatorDTO

	decoder := json.NewDecoder(ctx.Request.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&filter); err != nil {
		if a.debug_ {
			fmt.Println(err.Error())
		}
		utils.RespondError(ctx, http.StatusBadRequest, "invalid body", err.Error())
		return
	}

	if a.debug_ {
		fmt.Println(filter)
	}

	serviceData, err := a.repository_.FindService(filter.ServiceName)
	if err != nil {
		if a.debug_ {
			fmt.Println(err.Error())
		}
		utils.RespondError(ctx, http.StatusInternalServerError, "failed to find service", err.Error())
		return
	}

	if serviceData.ServiceID == 0 {
		if a.debug_ {
			fmt.Printf("service not found (service_name: %s)\n", filter.ServiceName)
		}
		a.repository_.CreateService(filter.ServiceName)
		utils.RespondError(ctx, http.StatusNotFound, "policy not found", "policy not found")
		return
	}

	user_id, err := a.repository_.FindUserIdByName(filter.UserName)
	if err != nil {
		if a.debug_ {
			fmt.Println(err.Error())
		}
		utils.RespondError(ctx, http.StatusInternalServerError, "failed to find user", err.Error())
		return
	}

	if user_id == nil {
		if a.debug_ {
			fmt.Println("internal aggregation error: nil ids")
		}
		utils.RespondError(ctx, http.StatusInternalServerError, "internal aggregation error: nil ids", "internal aggregation error: nil ids")
		return
	}

	policyData, err := a.repository_.FindPolicy(*user_id)
	if err != nil {
		if a.debug_ {
			fmt.Println(err.Error())
		}
		utils.RespondError(ctx, http.StatusInternalServerError, "failed to find policy", err.Error())
		return
	}

	if policyData.PolicyID == 0 {
		if a.debug_ {
			fmt.Printf("policy not found (user_id: %d)\n", *user_id)
		}
		utils.RespondError(ctx, http.StatusNotFound, "policy not found", "policy not found")
		return
	}

	ruleData, err := a.repository_.FindRule(serviceData.ServiceID, policyData.PolicyID)
	if err != nil {
		if a.debug_ {
			fmt.Println(err.Error())
		}
		utils.RespondError(ctx, http.StatusInternalServerError, "failed to find rule", err.Error())
		return
	}

	if ruleData == nil {
		if a.debug_ {
			fmt.Println("rule is empty")
		}
		utils.RespondError(ctx, http.StatusInternalServerError, "rule is empty", "rule is empty")
		return
	}

	var conditions []Condition
	if err := json.Unmarshal(ruleData.Condition, &conditions); err != nil {
		if a.debug_ {
			fmt.Println(err.Error())
		}
		utils.RespondError(ctx, http.StatusInternalServerError, "invalid rule condition format", err.Error())
		return
	}

	response := PolicyMatchResponse{
		Policy: PolicyResponse{
			ID:     policyData.PolicyID,
			Name:   policyData.PolicyName,
			UserID: *policyData.UserID,
			Status: policyData.Status,
		},
		Rule: Rule{
			Conditions: conditions,
			Effect:     ruleData.Effect,
		},
	}

	utils.RespondSuccess(ctx, http.StatusOK, nil, response)
}
