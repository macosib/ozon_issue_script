package ozonAPI

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type CampaignProduct struct {
	Client *OzonAPIClient
}

func NewCampaignProduct(Client *OzonAPIClient) *CampaignProduct {
	return &CampaignProduct{Client: Client}
}

type AddProductToCampaignRequest struct {
	CampaignID string `json:"campaign_id"`
	ProductID  string `json:"product_id"`
}

type Product struct {
	Sku        string      `json:"sku"`
	Bid        string      `json:"bid"`
	Categories []Category  `json:"categories,omitempty"`
	Title      string      `json:"title"`
	GroupID    string      `json:"groupId"`
	StopWords  []string    `json:"stopWords,omitempty"`
	Phrases    []PhraseBid `json:"phrases,omitempty"`
}

type PhraseBid struct {
	Bid             string `json:"bid"`
	Phrase          string `json:"phrase"`
	RelevanceStatus string `json:"relevanceStatus"`
}

type Category struct {
	Bid        uint64 `json:"bid"`
	CategoryId uint64 `json:"categoryIdid"`
}

type RequestBodyAddProduct struct {
	Bids []Product `json:"bids"`
}

// NewRequestBodyAddProduct Метод для создания тела запроса к методу добавления товаров в акции.
func NewRequestBodyAddProduct(productSkus []string) RequestBodyAddProduct {
	var bids []Product
	for _, sku := range productSkus {
		productBid := Product{
			Sku:       sku,
			Bid:       "30000000",
			GroupID:   "0",
			Phrases:   []PhraseBid{},
			StopWords: []string{},
		}
		bids = append(bids, productBid)
	}

	return RequestBodyAddProduct{
		Bids: bids,
	}
}

type RequestBodyRemoveProduct struct {
	Sku []string `json:"sku"`
}

// NewRequestBodyAddProduct Метод для создания тела запроса к методу добавления товаров в акции.
func NewRequestBodyRemoveProduct(productSkus []string) RequestBodyRemoveProduct {
	return RequestBodyRemoveProduct{
		Sku: productSkus,
	}
}

type ResponseCampaignProduct struct {
	Products []Product `json:"products"`
}

// Add Метод для добавления товаров в акции.
func (c *CampaignProduct) Add(campaignID string, products []string) error {
	endpoint := fmt.Sprintf("/campaign/%s/products", campaignID)
	payload := NewRequestBodyAddProduct(products)
	_, err := c.Client.sendRequest(http.MethodPost, endpoint, payload)
	if err != nil {
		return fmt.Errorf("ошибка при добавлении товаров в рекламную компанию: %w", err)
	}
	return nil
}

// Remove Метод метод для удаления товаров из акции.
func (c *CampaignProduct) Remove(campaignID string, products []string) error {
	endpoint := fmt.Sprintf("/campaign/%s/products/delete", campaignID)
	payload := NewRequestBodyRemoveProduct(products)
	_, err := c.Client.sendRequest(http.MethodPost, endpoint, payload)
	if err != nil {
		return fmt.Errorf("ошибка при удалении товаров из рекламной компании: %w", err)
	}
	return nil
}

// Get Метод для получения товаров из рекламной кампании.
func (c *CampaignProduct) Get(campaignID string) ([]Product, error) {

	var campaignProductResp ResponseCampaignProduct

	endpoint := fmt.Sprintf("/campaign/%s/products", campaignID)
	response, err := c.Client.sendRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("ошибка при запросе к API на получение списка товаров в рекламаной компании: %w", err)
	}

	err = json.NewDecoder(response.Body).Decode(&campaignProductResp)
	defer response.Body.Close()
	if err != nil {
		errMsg := fmt.Errorf(" ошибка при чтении ответа: %v", err)
		return nil, errMsg
	}
	return campaignProductResp.Products, nil
}

// GetCampaignProductSku Метод для получения sku товаров из рекламной кампании.
func (c *CampaignProduct) GetCampaignProductSku(campaignID string) ([]string, error) {

	var result []string

	products, err := c.Get(campaignID)

	if err != nil {
		errMsg := fmt.Errorf("ошибка при получении списка товаров: %v", err)
		return nil, errMsg
	}

	for _, v := range products {
		result = append(result, v.Sku)
	}

	return result, nil

}
