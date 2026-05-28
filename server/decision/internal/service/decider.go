package service

import (
	"decision/internal/client"
	"decision/internal/model"
	"fmt"
	"net"
	"regexp"
	"strings"
	"time"
)

type Decider struct {
	client          client.DecisionClient
	defaultDecision bool
}

func NewDecider(client client.DecisionClient, defaultDecision bool) *Decider {
	return &Decider{
		client:          client,
		defaultDecision: defaultDecision,
	}
}

func (d *Decider) Fallback(response *model.PolicyMatchResponse, request *model.ApiRequest) (*model.Decision, error) {
	requestId, err := d.client.CreateFallbackRequest(request)
	if err != nil {
		return nil, fmt.Errorf("error in CreateRequest: %s", err.Error())
	}

	decisionId, err := d.client.CreateFallbackDecision(requestId, d.defaultDecision)
	if err != nil {
		return nil, fmt.Errorf("error in CreateDecision: %s", err.Error())
	}
	return &model.Decision{
		Result: d.defaultDecision,
		Policy: struct {
			Id   uint
			Name string
		}{
			Id:   0,
			Name: "unknown",
		},
		RequestId:  requestId,
		DecisionId: decisionId,
	}, nil
}

func (d *Decider) Evaluate(request *model.ApiRequest) (*model.Decision, error) {
	response, err := d.client.GetRule(request)
	if err != nil {
		return d.Fallback(response, request)
	}

	type reqRes struct {
		id  uint
		err error
	}

	type ruleRes struct {
		ok  bool
		err error
	}

	reqCh := make(chan reqRes, 1)
	ruleCh := make(chan ruleRes, 1)

	go func() {
		id, err := d.client.CreateRequest(
			response.Policy.UserID,
			request,
		)

		reqCh <- reqRes{id: id, err: err}
	}()

	go func() {
		rule, err := model.ParseRule(response.Rule)
		if err != nil {
			ruleCh <- ruleRes{err: fmt.Errorf("ParseRule: %w", err)}
			return
		}

		ok, err := d.ApplyRule(request, rule)

		ruleCh <- ruleRes{ok: ok, err: err}
	}()

	req := <-reqCh
	if req.err != nil {
		return nil, fmt.Errorf("CreateRequest: %w", req.err)
	}

	rule := <-ruleCh
	if rule.err != nil {
		return nil, fmt.Errorf("ApplyRule: %w", rule.err)
	}

	decisionID, err := d.client.CreateDecision(
		req.id,
		response.Policy.ID,
		rule.ok,
	)
	if err != nil {
		return nil, fmt.Errorf("CreateDecision: %w", err)
	}

	return &model.Decision{
		Result: rule.ok,
		Policy: struct {
			Id   uint
			Name string
		}{
			Id:   response.Policy.ID,
			Name: response.Policy.Name,
		},
		RequestId:  req.id,
		DecisionId: decisionID,
	}, nil
}

func (d *Decider) ApplyRule(req *model.ApiRequest, rule *model.Rule) (bool, error) {

	for _, cond := range rule.Conditions {

		var res bool
		var err error

		switch cond.Type {

		case "gid":
			res, err = d.CheckGID(req, &cond)

		case "groups":
			res, err = d.CheckGroups(req, &cond)

		case "ip":
			res, err = d.CheckIP(req, &cond)

		case "hostname":
			res, err = d.CheckHostname(req, &cond)

		case "timestamp":
			res, err = d.CheckTimestamp(req, &cond)

		case "weekday":
			res, err = d.CheckWeekday(req, &cond)

		default:
			return false, fmt.Errorf("unknown condition type: %s", cond.Type)
		}

		if err != nil {
			return false, err
		}

		if !res {
			return !rule.Effect, nil
		}

	}

	return rule.Effect, nil
}

func (d *Decider) CheckGID(req *model.ApiRequest, cond *model.Condition) (bool, error) {
	val, ok := cond.Value.(float64)
	if !ok {
		return false, fmt.Errorf("invalid value type for gid")
	}

	gid := uint(val)

	switch cond.Operator {

	case model.OpEquals:
		return req.User.GID == gid, nil

	case model.OpNotEquals:
		return req.User.GID != gid, nil

	default:
		return false, fmt.Errorf("unknown operator: %s", cond.Operator)
	}
}

func (d *Decider) CheckGroups(req *model.ApiRequest, cond *model.Condition) (bool, error) {
	groups := req.User.Groups

	switch cond.Operator {

	case model.OpEquals:
		val, ok := cond.Value.(string)
		if !ok {
			return false, fmt.Errorf("invalid value for group")
		}
		for _, g := range groups {
			if g == val {
				return true, nil
			}
		}
		return false, nil

	case model.OpContains:
		val, ok := cond.Value.(string)
		if !ok {
			return false, fmt.Errorf("invalid value for group")
		}
		for _, g := range groups {
			if g == val {
				return true, nil
			}
		}
		return false, nil

	default:
		return false, fmt.Errorf("unsupported operator for group: %s", cond.Operator)
	}
}

func (d *Decider) CheckIP(req *model.ApiRequest, cond *model.Condition) (bool, error) {
	ip := net.ParseIP(req.Host.IP)
	if ip == nil {
		return false, fmt.Errorf("invalid request ip")
	}

	switch cond.Operator {

	case model.OpEquals:
		val, ok := cond.Value.(string)
		if !ok {
			return false, fmt.Errorf("invalid value for ip")
		}
		return req.Host.IP == val, nil

	case model.OpIn:
		val, ok := cond.Value.(string)
		if !ok {
			return false, fmt.Errorf("invalid cidr")
		}

		_, ipNet, err := net.ParseCIDR(val)
		if err != nil {
			return false, err
		}

		return ipNet.Contains(ip), nil

	case model.OpNotIn:
		val, ok := cond.Value.(string)
		if !ok {
			return false, fmt.Errorf("invalid cidr")
		}

		_, ipNet, err := net.ParseCIDR(val)
		if err != nil {
			return false, err
		}

		return !ipNet.Contains(ip), nil

	default:
		return false, fmt.Errorf("unsupported operator for ip: %s", cond.Operator)
	}
}

func (d *Decider) CheckHostname(req *model.ApiRequest, cond *model.Condition) (bool, error) {

	val, ok := cond.Value.(string)
	if !ok {
		return false, fmt.Errorf("invalid hostname value")
	}

	switch cond.Operator {

	case model.OpEquals:
		return req.Host.HostName == val, nil

	case model.OpNotEquals:
		return req.Host.HostName != val, nil

	case model.OpRegex:
		return regexp.MatchString(val, req.Host.HostName)

	default:
		return false, fmt.Errorf("unsupported operator for hostname: %s", cond.Operator)
	}
}

func (d *Decider) CheckTimestamp(req *model.ApiRequest, cond *model.Condition) (bool, error) {
	val, ok := cond.Value.(string)
	if !ok {
		return false, fmt.Errorf("invalid time range")
	}

	parts := strings.Split(val, "-")
	if len(parts) != 2 {
		return false, fmt.Errorf("invalid time format, expected HH:MM-HH:MM")
	}

	start, err := time.Parse("15:04", parts[0])
	if err != nil {
		return false, err
	}

	end, err := time.Parse("15:04", parts[1])
	if err != nil {
		return false, err
	}

	// Используем время из req.Time.Timestamp как есть, без перевода в другую timezone.
	now := req.Time.Timestamp

	// Сравниваем только часы и минуты, сохраняя ту же Location,
	// которая уже присутствует в переданном времени.
	loc := now.Location()

	current := time.Date(
		now.Year(), now.Month(), now.Day(),
		now.Hour(), now.Minute(), 0, 0, loc,
	)

	startTime := time.Date(
		now.Year(), now.Month(), now.Day(),
		start.Hour(), start.Minute(), 0, 0, loc,
	)

	endTime := time.Date(
		now.Year(), now.Month(), now.Day(),
		end.Hour(), end.Minute(), 0, 0, loc,
	)

	switch cond.Operator {
	case model.OpBetween:
		return (current.Equal(startTime) || current.After(startTime)) &&
			(current.Equal(endTime) || current.Before(endTime)), nil

	default:
		return false, fmt.Errorf("unsupported operator for timestamp: %s", cond.Operator)
	}
}

func (d *Decider) CheckWeekday(req *model.ApiRequest, cond *model.Condition) (bool, error) {

	weekday := strings.ToLower(req.Time.Weekday)

	switch cond.Operator {

	case model.OpEquals:
		val, ok := cond.Value.(string)
		if !ok {
			return false, fmt.Errorf("invalid weekday")
		}

		return weekday == strings.ToLower(val), nil

	case model.OpIn:
		arr, ok := cond.Value.([]interface{})
		if !ok {
			return false, fmt.Errorf("invalid weekday list")
		}

		for _, v := range arr {
			s, ok := v.(string)
			if !ok {
				continue
			}

			if weekday == strings.ToLower(s) {
				return true, nil
			}
		}

		return false, nil

	default:
		return false, fmt.Errorf("unsupported operator for weekday: %s", cond.Operator)
	}
}

func (d *Decider) CheckService(req *model.ApiRequest, cond *model.Condition) (bool, error) {

	service := req.Service

	switch cond.Operator {

	case model.OpEquals:
		val, ok := cond.Value.(string)
		if !ok {
			return false, fmt.Errorf("invalid service")
		}

		return service == strings.ToLower(val), nil

	case model.OpIn:
		arr, ok := cond.Value.([]interface{})
		if !ok {
			return false, fmt.Errorf("invalid service list")
		}

		for _, v := range arr {
			s, ok := v.(string)
			if !ok {
				continue
			}

			if service == strings.ToLower(s) {
				return true, nil
			}
		}

		return false, nil

	default:
		return false, fmt.Errorf("unsupported operator for service: %s", cond.Operator)
	}
}
