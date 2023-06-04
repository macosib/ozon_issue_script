package ozonAPI

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type OzonAPIClient struct {
	BaseUrl           string
	AccessToken       string
	ApiIdPerformance  string
	apiKeyPerformance string
}

type RequestBodyAccessToken struct {
	ClientId     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	GrantType    string `json:"grant_type"`
}

type ResponseAccessToken struct {
	Access_token string `json:"access_token"`
	Expires_in   int64  `json:"expires_in"`
	Token_type   string `json:"token_type"`
}

func NewRequestBodyAccessToken(apiIdPerformance string, apiKeyPerformance string) *RequestBodyAccessToken {
	return &RequestBodyAccessToken{
		ClientId:     apiIdPerformance,
		ClientSecret: apiKeyPerformance,
		GrantType:    "client_credentials",
	}
}

func NewOzonAPIClient(apiIdPerformance string, apiKeyPerformance string) *OzonAPIClient {
	return &OzonAPIClient{
		BaseUrl:           "https://performance.ozon.ru/api/client",
		ApiIdPerformance:  apiIdPerformance,
		apiKeyPerformance: apiKeyPerformance,
		AccessToken:       "",
	}
}

func (c *OzonAPIClient) sendRequest(method, endpoint string, payload interface{}) (*http.Response, error) {
	requestBody := make([]byte, 0)

	if payload != nil {
		var err error
		requestBody, err = json.Marshal(payload)
		if err != nil {
			return nil, fmt.Errorf("ошибка при сериализации запроса: %w", err)
		}
	}
	url := c.BaseUrl + endpoint
	request, err := http.NewRequest(method, url, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("ошибка при создании запроса: %w", err)
	}

	request.Header.Set("Authorization", fmt.Sprintf("Bearer %v", c.AccessToken))
	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)

	if err != nil {
		return nil, fmt.Errorf("ошибка при отправке запроса: %w", err)
	}
	err = c.checkStatusCode(response)
	if err != nil {
		return nil, err
	}
	return response, nil
}

// SetAccessToken Получаем токен и записываем в структуру для запросов к API Ozon Performance.
func (c *OzonAPIClient) SetAccessToken() error {
	token, err := c.getAccessToken()
	if err != nil {
		log.Fatal(err)
	}
	c.AccessToken = token
	return nil
}

// getAccessToken Получаем токен для запросов к API Ozon Performance.
func (c *OzonAPIClient) getAccessToken() (string, error) {
	var result ResponseAccessToken

	endpoint := "/token"

	requestBody := NewRequestBodyAccessToken(c.ApiIdPerformance, c.apiKeyPerformance)
	response, err := c.sendRequest(http.MethodPost, endpoint, requestBody)
	if err != nil {
		errMessage := fmt.Errorf("ошибка при отправке запроса к API, client: %s,  message: %s", c.ApiIdPerformance, err)
		return "", errMessage
	}

	err = json.NewDecoder(response.Body).Decode(&result)
	defer response.Body.Close()
	if err != nil {
		errMessage := fmt.Errorf("не удалось распознать ответ от Performance Ozon API, client: %s  status code: %d, message: %s", c.ApiIdPerformance, response.StatusCode, string(err.Error()))
		return "", errMessage
	}

	return result.Access_token, nil
}

// checkStatusCode Проверяем статус код запроса.
func (c *OzonAPIClient) checkStatusCode(response *http.Response) error {
	if response.StatusCode != http.StatusOK {
		defer response.Body.Close()
		reader := bufio.NewReader(response.Body)
		err, _ := io.ReadAll(reader)
		errMessage := fmt.Errorf("не удалось получить ответ от Performance Ozon API, client: %s  status code: %d, message: %s", c.ApiIdPerformance, response.StatusCode, err)
		return errMessage
	}
	return nil
}
