package service

import (
	"decision/internal/model"
	"decision/internal/service"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type fakeDecisionClient struct {
	getRuleResp               *model.PolicyMatchResponse
	getRuleErr                error
	createRequestID           uint
	createRequestErr          error
	createDecisionID          uint
	createDecisionErr         error
	createFallbackRequestID   uint
	createFallbackRequestErr  error
	createFallbackDecisionID  uint
	createFallbackDecisionErr error
	createdUserID             uint
	createdRequestID          uint
	createdPolicyID           uint
	createdResult             bool
}

func (f *fakeDecisionClient) GetRule(request *model.ApiRequest) (*model.PolicyMatchResponse, error) {
	return f.getRuleResp, f.getRuleErr
}

func (f *fakeDecisionClient) CreateRequest(userID uint, request *model.ApiRequest) (uint, error) {
	f.createdUserID = userID
	return f.createRequestID, f.createRequestErr
}

func (f *fakeDecisionClient) CreateDecision(requestID uint, policyID uint, result bool) (uint, error) {
	f.createdRequestID = requestID
	f.createdPolicyID = policyID
	f.createdResult = result
	return f.createDecisionID, f.createDecisionErr
}

func (f *fakeDecisionClient) CreateFallbackRequest(request *model.ApiRequest) (uint, error) {
	return f.createFallbackRequestID, f.createFallbackRequestErr
}

func (f *fakeDecisionClient) CreateFallbackDecision(requestID uint, result bool) (uint, error) {
	f.createdRequestID = requestID
	f.createdResult = result
	return f.createFallbackDecisionID, f.createFallbackDecisionErr
}

func buildRequest() *model.ApiRequest {
	req := &model.ApiRequest{}

	req.User.Name = "roman"
	req.User.UID = 1001
	req.User.GID = 1001
	req.User.Groups = []string{
		"sudo",
		"dev",
	}

	req.Host.IP = "192.168.1.10"
	req.Host.HostName = "prod-1"

	req.Service = "sshd"

	req.Time.Timestamp = time.Date(
		2026,
		time.April,
		15,
		12,
		0,
		0,
		0,
		time.UTC,
	)

	req.Time.Weekday = "tue"

	return req
}

func buildPolicyResponse(rule model.Rule) *model.PolicyMatchResponse {
	return &model.PolicyMatchResponse{
		Policy: model.PolicyResponse{
			ID:     22,
			Name:   "ssh-prod",
			UserID: 77,
			Status: true,
		},
		Rule: rule,
	}
}

func TestNewDecider(t *testing.T) {
	client := &fakeDecisionClient{}

	decider := service.NewDecider(client, true)

	assert.NotNil(t, decider)
}

func TestFallback(t *testing.T) {
	req := buildRequest()

	t.Run("success uses default decision and unknown policy", func(t *testing.T) {
		client := &fakeDecisionClient{
			createFallbackRequestID:  10,
			createFallbackDecisionID: 20,
		}
		decider := service.NewDecider(client, true)

		got, err := decider.Fallback(nil, req)

		assert.NoError(t, err)
		assert.Equal(t, true, got.Result)
		assert.Equal(t, uint(0), got.Policy.Id)
		assert.Equal(t, "unknown", got.Policy.Name)
		assert.Equal(t, uint(10), got.RequestId)
		assert.Equal(t, uint(20), got.DecisionId)
		assert.Equal(t, uint(10), client.createdRequestID)
		assert.Equal(t, true, client.createdResult)
	})

	t.Run("create fallback request error", func(t *testing.T) {
		decider := service.NewDecider(&fakeDecisionClient{
			createFallbackRequestErr: errors.New("request failed"),
		}, false)

		got, err := decider.Fallback(nil, req)

		assert.Nil(t, got)
		assert.ErrorContains(t, err, "error in CreateRequest")
	})

	t.Run("create fallback decision error", func(t *testing.T) {
		decider := service.NewDecider(&fakeDecisionClient{
			createFallbackRequestID:   10,
			createFallbackDecisionErr: errors.New("decision failed"),
		}, false)

		got, err := decider.Fallback(nil, req)

		assert.Nil(t, got)
		assert.ErrorContains(t, err, "error in CreateDecision")
	})
}

func TestEvaluate(t *testing.T) {
	req := buildRequest()
	allowRule := model.Rule{
		Effect: true,
		Conditions: []model.Condition{
			{
				Type:     "gid",
				Operator: model.OpEquals,
				Value:    float64(1001),
			},
		},
	}

	t.Run("success creates request and decision", func(t *testing.T) {
		client := &fakeDecisionClient{
			getRuleResp:      buildPolicyResponse(allowRule),
			createRequestID:  11,
			createDecisionID: 33,
		}
		decider := service.NewDecider(client, false)

		got, err := decider.Evaluate(req)

		assert.NoError(t, err)
		assert.Equal(t, true, got.Result)
		assert.Equal(t, uint(22), got.Policy.Id)
		assert.Equal(t, "ssh-prod", got.Policy.Name)
		assert.Equal(t, uint(11), got.RequestId)
		assert.Equal(t, uint(33), got.DecisionId)
		assert.Equal(t, uint(77), client.createdUserID)
		assert.Equal(t, uint(11), client.createdRequestID)
		assert.Equal(t, uint(22), client.createdPolicyID)
		assert.Equal(t, true, client.createdResult)
	})

	t.Run("get rule error falls back", func(t *testing.T) {
		client := &fakeDecisionClient{
			getRuleErr:               errors.New("aggregator down"),
			createFallbackRequestID:  44,
			createFallbackDecisionID: 55,
		}
		decider := service.NewDecider(client, false)

		got, err := decider.Evaluate(req)

		assert.NoError(t, err)
		assert.Equal(t, false, got.Result)
		assert.Equal(t, uint(44), got.RequestId)
		assert.Equal(t, uint(55), got.DecisionId)
	})

	t.Run("create request error", func(t *testing.T) {
		decider := service.NewDecider(&fakeDecisionClient{
			getRuleResp:      buildPolicyResponse(allowRule),
			createRequestErr: errors.New("request failed"),
		}, false)

		got, err := decider.Evaluate(req)

		assert.Nil(t, got)
		assert.ErrorContains(t, err, "CreateRequest")
	})

	t.Run("parse rule error", func(t *testing.T) {
		badRule := model.Rule{
			Effect: true,
			Conditions: []model.Condition{
				{
					Type:     "gid",
					Operator: model.OpEquals,
					Value:    func() {},
				},
			},
		}
		decider := service.NewDecider(&fakeDecisionClient{
			getRuleResp:     buildPolicyResponse(badRule),
			createRequestID: 11,
		}, false)

		got, err := decider.Evaluate(req)

		assert.Nil(t, got)
		assert.ErrorContains(t, err, "ApplyRule")
		assert.ErrorContains(t, err, "ParseRule")
	})

	t.Run("apply rule error", func(t *testing.T) {
		badRule := model.Rule{
			Effect: true,
			Conditions: []model.Condition{
				{
					Type:     "gid",
					Operator: "bad",
					Value:    float64(1001),
				},
			},
		}
		decider := service.NewDecider(&fakeDecisionClient{
			getRuleResp:     buildPolicyResponse(badRule),
			createRequestID: 11,
		}, false)

		got, err := decider.Evaluate(req)

		assert.Nil(t, got)
		assert.ErrorContains(t, err, "ApplyRule")
	})

	t.Run("create decision error", func(t *testing.T) {
		decider := service.NewDecider(&fakeDecisionClient{
			getRuleResp:       buildPolicyResponse(allowRule),
			createRequestID:   11,
			createDecisionErr: errors.New("decision failed"),
		}, false)

		got, err := decider.Evaluate(req)

		assert.Nil(t, got)
		assert.ErrorContains(t, err, "CreateDecision")
	})
}

func TestApplyRule(t *testing.T) {
	d := &service.Decider{}
	req := buildRequest()

	tests := []struct {
		name    string
		rule    *model.Rule
		want    bool
		wantErr bool
	}{
		{
			name: "all matched allow",
			rule: &model.Rule{
				Effect: true,
				Conditions: []model.Condition{
					{
						Type:     "gid",
						Operator: model.OpEquals,
						Value:    float64(1001),
					},
					{
						Type:     "groups",
						Operator: model.OpContains,
						Value:    "sudo",
					},
					{
						Type:     "ip",
						Operator: model.OpIn,
						Value:    "192.168.1.0/24",
					},
				},
			},
			want: true,
		},
		{
			name: "inverse effect",
			rule: &model.Rule{
				Effect: true,
				Conditions: []model.Condition{
					{
						Type:     "gid",
						Operator: model.OpEquals,
						Value:    float64(9999),
					},
				},
			},
			want: false,
		},
		{
			name: "effect false",
			rule: &model.Rule{
				Effect: false,
				Conditions: []model.Condition{
					{
						Type:     "hostname",
						Operator: model.OpEquals,
						Value:    "prod-1",
					},
				},
			},
			want: false,
		},
		{
			name: "unknown condition",
			rule: &model.Rule{
				Effect: true,
				Conditions: []model.Condition{
					{
						Type:     "unknown",
						Operator: model.OpEquals,
						Value:    "x",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "condition error",
			rule: &model.Rule{
				Effect: true,
				Conditions: []model.Condition{
					{
						Type:     "ip",
						Operator: model.OpIn,
						Value:    "bad-cidr",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "early return",
			rule: &model.Rule{
				Effect: true,
				Conditions: []model.Condition{
					{
						Type:     "gid",
						Operator: model.OpEquals,
						Value:    float64(9999),
					},
					{
						Type:     "ip",
						Operator: model.OpEquals,
						Value:    "192.168.1.10",
					},
				},
			},
			want: false,
		},
		{
			name: "multiple conditions",
			rule: &model.Rule{
				Effect: true,
				Conditions: []model.Condition{
					{
						Type:     "gid",
						Operator: model.OpEquals,
						Value:    float64(1001),
					},
					{
						Type:     "groups",
						Operator: model.OpContains,
						Value:    "sudo",
					},
					{
						Type:     "hostname",
						Operator: model.OpRegex,
						Value:    "^prod",
					},
					{
						Type:     "timestamp",
						Operator: model.OpBetween,
						Value:    "09:00-18:00",
					},
					{
						Type:     "weekday",
						Operator: model.OpEquals,
						Value:    "tue",
					},
				},
			},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := d.ApplyRule(req, tt.rule)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestCheckGID(t *testing.T) {
	d := &service.Decider{}
	req := buildRequest()

	tests := []struct {
		name    string
		cond    model.Condition
		want    bool
		wantErr bool
	}{
		{
			name: "equals",
			cond: model.Condition{
				Operator: model.OpEquals,
				Value:    float64(1001),
			},
			want: true,
		},
		{
			name: "not equals",
			cond: model.Condition{
				Operator: model.OpNotEquals,
				Value:    float64(2000),
			},
			want: true,
		},
		{
			name: "invalid type",
			cond: model.Condition{
				Operator: model.OpEquals,
				Value:    "1001",
			},
			wantErr: true,
		},
		{
			name: "unsupported operator",
			cond: model.Condition{
				Operator: "bad",
				Value:    float64(1001),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := d.CheckGID(req, &tt.cond)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestCheckGroups(t *testing.T) {
	d := &service.Decider{}
	req := buildRequest()

	tests := []struct {
		name    string
		cond    model.Condition
		want    bool
		wantErr bool
	}{
		{
			name: "contains",
			cond: model.Condition{
				Operator: model.OpContains,
				Value:    "sudo",
			},
			want: true,
		},
		{
			name: "equals",
			cond: model.Condition{
				Operator: model.OpEquals,
				Value:    "dev",
			},
			want: true,
		},
		{
			name: "not found",
			cond: model.Condition{
				Operator: model.OpContains,
				Value:    "admin",
			},
			want: false,
		},
		{
			name: "equals not found",
			cond: model.Condition{
				Operator: model.OpEquals,
				Value:    "admin",
			},
			want: false,
		},
		{
			name: "invalid value",
			cond: model.Condition{
				Operator: model.OpContains,
				Value:    123,
			},
			wantErr: true,
		},
		{
			name: "invalid equals value",
			cond: model.Condition{
				Operator: model.OpEquals,
				Value:    123,
			},
			wantErr: true,
		},
		{
			name: "unsupported operator",
			cond: model.Condition{
				Operator: "bad",
				Value:    "sudo",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := d.CheckGroups(req, &tt.cond)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestCheckIP(t *testing.T) {
	d := &service.Decider{}

	t.Run("invalid request ip", func(t *testing.T) {
		req := buildRequest()
		req.Host.IP = "bad-ip"

		_, err := d.CheckIP(req, &model.Condition{
			Operator: model.OpEquals,
			Value:    "1.1.1.1",
		})

		assert.Error(t, err)
	})

	tests := []struct {
		name    string
		cond    model.Condition
		want    bool
		wantErr bool
	}{
		{
			name: "equals",
			cond: model.Condition{
				Operator: model.OpEquals,
				Value:    "192.168.1.10",
			},
			want: true,
		},
		{
			name: "in cidr",
			cond: model.Condition{
				Operator: model.OpIn,
				Value:    "192.168.1.0/24",
			},
			want: true,
		},
		{
			name: "not in cidr",
			cond: model.Condition{
				Operator: model.OpNotIn,
				Value:    "10.0.0.0/8",
			},
			want: true,
		},
		{
			name: "invalid equals value",
			cond: model.Condition{
				Operator: model.OpEquals,
				Value:    123,
			},
			wantErr: true,
		},
		{
			name: "invalid cidr type",
			cond: model.Condition{
				Operator: model.OpIn,
				Value:    123,
			},
			wantErr: true,
		},
		{
			name: "invalid not in cidr type",
			cond: model.Condition{
				Operator: model.OpNotIn,
				Value:    123,
			},
			wantErr: true,
		},
		{
			name: "invalid cidr",
			cond: model.Condition{
				Operator: model.OpIn,
				Value:    "bad",
			},
			wantErr: true,
		},
		{
			name: "invalid not in cidr",
			cond: model.Condition{
				Operator: model.OpNotIn,
				Value:    "bad",
			},
			wantErr: true,
		},
		{
			name: "unsupported operator",
			cond: model.Condition{
				Operator: "bad",
				Value:    "192.168.1.1",
			},
			wantErr: true,
		},
	}

	req := buildRequest()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := d.CheckIP(req, &tt.cond)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestCheckHostname(t *testing.T) {
	d := &service.Decider{}
	req := buildRequest()

	tests := []struct {
		name    string
		cond    model.Condition
		want    bool
		wantErr bool
	}{
		{
			name: "equals",
			cond: model.Condition{
				Operator: model.OpEquals,
				Value:    "prod-1",
			},
			want: true,
		},
		{
			name: "not equals",
			cond: model.Condition{
				Operator: model.OpNotEquals,
				Value:    "dev-1",
			},
			want: true,
		},
		{
			name: "regex",
			cond: model.Condition{
				Operator: model.OpRegex,
				Value:    "^prod-.*",
			},
			want: true,
		},
		{
			name: "regex no match",
			cond: model.Condition{
				Operator: model.OpRegex,
				Value:    "^dev",
			},
			want: false,
		},
		{
			name: "invalid regex",
			cond: model.Condition{
				Operator: model.OpRegex,
				Value:    "[",
			},
			wantErr: true,
		},
		{
			name: "invalid value",
			cond: model.Condition{
				Operator: model.OpEquals,
				Value:    123,
			},
			wantErr: true,
		},
		{
			name: "unsupported operator",
			cond: model.Condition{
				Operator: "bad",
				Value:    "prod-1",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := d.CheckHostname(req, &tt.cond)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestCheckTimestamp(t *testing.T) {
	d := &service.Decider{}
	req := buildRequest()

	tests := []struct {
		name    string
		cond    model.Condition
		want    bool
		wantErr bool
	}{
		{
			name: "between",
			cond: model.Condition{
				Operator: model.OpBetween,
				Value:    "09:00-18:00",
			},
			want: true,
		},
		{
			name: "outside",
			cond: model.Condition{
				Operator: model.OpBetween,
				Value:    "18:01-20:00",
			},
			want: false,
		},
		{
			name: "boundary start",
			cond: model.Condition{
				Operator: model.OpBetween,
				Value:    "12:00-18:00",
			},
			want: true,
		},
		{
			name: "boundary end",
			cond: model.Condition{
				Operator: model.OpBetween,
				Value:    "09:00-12:00",
			},
			want: true,
		},
		{
			name: "invalid format",
			cond: model.Condition{
				Operator: model.OpBetween,
				Value:    "bad",
			},
			wantErr: true,
		},
		{
			name: "invalid start",
			cond: model.Condition{
				Operator: model.OpBetween,
				Value:    "99:00-18:00",
			},
			wantErr: true,
		},
		{
			name: "invalid end",
			cond: model.Condition{
				Operator: model.OpBetween,
				Value:    "09:00-99:00",
			},
			wantErr: true,
		},
		{
			name: "invalid type",
			cond: model.Condition{
				Operator: model.OpBetween,
				Value:    123,
			},
			wantErr: true,
		},
		{
			name: "unsupported operator",
			cond: model.Condition{
				Operator: "bad",
				Value:    "09:00-18:00",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := d.CheckTimestamp(req, &tt.cond)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestCheckWeekday(t *testing.T) {
	d := &service.Decider{}
	req := buildRequest()

	tests := []struct {
		name    string
		cond    model.Condition
		want    bool
		wantErr bool
	}{
		{
			name: "equals",
			cond: model.Condition{
				Operator: model.OpEquals,
				Value:    "tue",
			},
			want: true,
		},
		{
			name: "case insensitive",
			cond: model.Condition{
				Operator: model.OpEquals,
				Value:    "TUE",
			},
			want: true,
		},
		{
			name: "in list",
			cond: model.Condition{
				Operator: model.OpIn,
				Value: []interface{}{
					"mon",
					"tue",
				},
			},
			want: true,
		},
		{
			name: "not in list",
			cond: model.Condition{
				Operator: model.OpIn,
				Value: []interface{}{
					"sat",
					"sun",
				},
			},
			want: false,
		},
		{
			name: "invalid equals value",
			cond: model.Condition{
				Operator: model.OpEquals,
				Value:    123,
			},
			wantErr: true,
		},
		{
			name: "invalid list",
			cond: model.Condition{
				Operator: model.OpIn,
				Value:    "tue",
			},
			wantErr: true,
		},
		{
			name: "invalid list items",
			cond: model.Condition{
				Operator: model.OpIn,
				Value: []interface{}{
					1,
					2,
					"wed",
				},
			},
			want: false,
		},
		{
			name: "unsupported operator",
			cond: model.Condition{
				Operator: "bad",
				Value:    "tue",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := d.CheckWeekday(req, &tt.cond)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestCheckService(t *testing.T) {
	d := &service.Decider{}
	req := buildRequest()

	tests := []struct {
		name    string
		cond    model.Condition
		want    bool
		wantErr bool
	}{
		{
			name: "equals",
			cond: model.Condition{
				Operator: model.OpEquals,
				Value:    "sshd",
			},
			want: true,
		},
		{
			name: "in list",
			cond: model.Condition{
				Operator: model.OpIn,
				Value: []interface{}{
					"nginx",
					"sshd",
				},
			},
			want: true,
		},
		{
			name: "not found",
			cond: model.Condition{
				Operator: model.OpIn,
				Value: []interface{}{
					"nginx",
					"apache",
				},
			},
			want: false,
		},
		{
			name: "invalid list items",
			cond: model.Condition{
				Operator: model.OpIn,
				Value: []interface{}{
					123,
					"api",
				},
			},
			want: false,
		},
		{
			name: "invalid equals value",
			cond: model.Condition{
				Operator: model.OpEquals,
				Value:    123,
			},
			wantErr: true,
		},
		{
			name: "invalid list",
			cond: model.Condition{
				Operator: model.OpIn,
				Value:    "sshd",
			},
			wantErr: true,
		},
		{
			name: "unsupported operator",
			cond: model.Condition{
				Operator: "bad",
				Value:    "sshd",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := d.CheckService(req, &tt.cond)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
