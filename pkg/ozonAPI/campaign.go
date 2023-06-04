package ozonAPI

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Campaign struct {
	Client *OzonAPIClient
}

func NewCampaign(Client *OzonAPIClient) *Campaign {
	return &Campaign{Client: Client}
}

type RequestCampaignCreate struct {
	Title               string `json:"title,omitempty"`
	FromDate            string `json:"fromDate,omitempty"`
	ToDate              string `json:"toDate,omitempty"`
	DailyBudget         string `json:"dailyBudget"`
	ExpenseStrategy     string `json:"expenseStrategy,omitempty"`
	Placement           string `json:"placement,omitempty"`
	ProductAutopilot    string `json:"productAutopilotStrategy,omitempty"`
	ProductCampaignMode string `json:"productCampaignMode,omitempty"`
}

type ResponseCreateCampaign struct {
	CampaignId string `json:"campaignId"`
}

// NewRequestBodyCreateCampaign Метод для создания тела запроса к методу создания компании.
func NewRequestBodyCreateCampaign(title string) *RequestCampaignCreate {
	return &RequestCampaignCreate{
		Title:               title,
		Placement:           "PLACEMENT_SEARCH_AND_CATEGORY",
		DailyBudget:         "550000000",
		ExpenseStrategy:     "DAILY_BUDGET",
		ProductCampaignMode: "PRODUCT_CAMPAIGN_MODE_AUTO",
	}
}

func getCreateCampaignTitle(campaignType string, isAdditionalCampaign bool) string {
	title := ""
	if campaignType == "CPC" {
		title = "Рекламная кампания с оплатой за клики - CPC. API"
	} else {
		title = "Рекламная кампания с оплатой за показы - CPM. API"
	}

	if !isAdditionalCampaign {
		return fmt.Sprintf("%v Main", title)
	}
	return fmt.Sprintf("%v Additional", title)
}

func getCreateCampaignEndpoint(campaignType string) string {
	if campaignType == "CPC" {
		return "/campaign/cpc/product"
	}
	return "/campaign/cpm/product"
}

// Create Метод создания рекламной компании.
// Рекламная кампания с оплатой за показы - CPM
// Рекламная кампания с оплатой за клики - CPC
func (c *Campaign) Create(campaignType string, isAdditionalCampaign bool) (string, error) {

	var createCampaignResponse ResponseCreateCampaign

	title := getCreateCampaignTitle(campaignType, isAdditionalCampaign)
	payload := NewRequestBodyCreateCampaign(title)
	endpoint := getCreateCampaignEndpoint(campaignType)
	response, err := c.Client.sendRequest(http.MethodPost, endpoint, payload)

	if err != nil {
		return "", fmt.Errorf("ошибка при запросе к API на создание рекламной компании: %w", err)
	}

	err = json.NewDecoder(response.Body).Decode(&createCampaignResponse)
	defer response.Body.Close()
	if err != nil {
		return "", fmt.Errorf("ошибка при десериализации ответа от API: %w", err)
	}

	return createCampaignResponse.CampaignId, nil
}

// Activate Метод активации рекламной компании.
func (c *Campaign) Activate(campaignId string) error {
	_, err := c.Client.sendRequest(http.MethodPost, fmt.Sprintf("/campaigns/%s/activate", campaignId), nil)
	if err != nil {
		return fmt.Errorf("failed to activate campaign: %w", err)
	}

	return nil
}

// Deactivate Метод деактивации рекламной компании.
func (c *Campaign) Deactivate(campaignId string) error {
	_, err := c.Client.sendRequest(http.MethodPost, fmt.Sprintf("/campaigns/%s/deactivate", campaignId), nil)
	if err != nil {
		return fmt.Errorf("failed to deactivate campaign: %w", err)
	}

	return nil
}
