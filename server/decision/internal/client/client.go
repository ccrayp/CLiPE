package client

import (
	"bytes"
	"decision/internal/model"
	"decision/pkg/config"
	"decision/pkg/utils"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type createDecisionDTO struct {
	RequestID uint `json:"request_id" binding:"required"`
	PolicyID  uint `json:"policy_id" binding:"required"`
	Result    bool `json:"result"`
}

type createRequestDTO struct {
	UserID  uint        `json:"user_id" binding:"required"`
	Context interface{} `json:"context"`
}

type Client struct {
	apiUrl     string
	httpClient http.Client
}

func NewClient(cfg *config.Config) *Client {
	return &Client{
		apiUrl:     "http://crud_server:8080" + "/api/v" + cfg.Server.ApiVersion + "/internal",
		httpClient: http.Client{},
	}
}

func (c *Client) CheckApiUrl() string {
	return c.apiUrl
}

func (c *Client) CreateRequest(userId uint, request *model.ApiRequest) (uint, error) {
	dto := createRequestDTO{
		UserID:  userId,
		Context: request,
	}

	var resp struct {
		utils.Common
		Data struct {
			ID uint `json:"id"`
		} `json:"data"`
	}

	err := c.doPost("/requests", dto, &resp)
	if err != nil {
		return 0, err
	}

	if !resp.Success {
		return 0, fmt.Errorf("request failed: %+v", resp)
	}

	return resp.Data.ID, nil
}

func (c *Client) CreateDecision(requestId uint, policyId uint, result bool) (uint, error) {
	dto := createDecisionDTO{
		RequestID: requestId,
		PolicyID:  policyId,
		Result:    result,
	}

	var resp struct {
		utils.Common
		Data struct {
			ID uint `json:"id"`
		} `json:"data"`
	}

	err := c.doPost("/decisions", dto, &resp)
	if err != nil {
		return 0, err
	}

	if !resp.Success {
		return 0, fmt.Errorf("request failed: %+v", resp)
	}

	return resp.Data.ID, nil
}

func (c *Client) CreateFallbackRequest(request *model.ApiRequest) (uint, error) {
	dto := createRequestDTO{
		Context: request,
	}

	var resp struct {
		utils.Common
		Data struct {
			ID uint `json:"id"`
		} `json:"data"`
	}

	err := c.doPost("/requests", dto, &resp)
	if err != nil {
		return 0, err
	}

	if !resp.Success {
		return 0, fmt.Errorf("request failed: %+v", resp)
	}

	return resp.Data.ID, nil
}

func (c *Client) CreateFallbackDecision(requestId uint, result bool) (uint, error) {
	dto := createDecisionDTO{
		RequestID: requestId,
		Result:    result,
	}

	var resp struct {
		utils.Common
		Data struct {
			ID uint `json:"id"`
		} `json:"data"`
	}

	err := c.doPost("/decisions", dto, &resp)
	if err != nil {
		return 0, err
	}

	if !resp.Success {
		return 0, fmt.Errorf("request failed: %+v", resp)
	}

	return resp.Data.ID, nil
}

func (c *Client) doPost(path string, dto any, out any) error {
	jsonData, err := json.Marshal(dto)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, c.apiUrl+path, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Internal-Token", os.Getenv("DECISION_TOKEN"))
	req.Header.Set("X-Caller", os.Getenv("DECISION_ID"))

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		var errResp utils.Error
		if err := json.Unmarshal(body, &errResp); err != nil {
			return fmt.Errorf("status: %d, body: %s", resp.StatusCode, string(body))
		}
		return fmt.Errorf("%s", errResp.Error)
	}

	return json.Unmarshal(body, out)
}

func (c *Client) GetRule(request *model.ApiRequest) (*model.PolicyMatchResponse, error) {
	dto := struct {
		UserName string `json:"user_name"`
	}{
		UserName: request.User.Name,
	}

	var resp struct {
		utils.Common
		Data model.PolicyMatchResponse `json:"data"`
	}

	err := c.doPost("/aggregator", dto, &resp)
	if err != nil {
		return nil, err
	}

	if !resp.Success {
		return nil, fmt.Errorf("request failed: %+v", resp)
	}

	return &resp.Data, nil
}
