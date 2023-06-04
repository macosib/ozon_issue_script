package main

import (
	"fmt"
	"log"
	"ozon_issue_script/pkg/ozonAPI"
)

func main() {

	apiIdPerformance := ""
	ApiKeyPerformance := ""

	// Получаем доступ к сервису работы с API Ozon Performance
	apiService := ozonAPI.NewOzonAPIClient(apiIdPerformance, ApiKeyPerformance)

	// Получаем AccessToken и сохраняем в инстанс OzonAPIClient
	err := apiService.SetAccessToken()
	if err != nil {
		errMsg := fmt.Errorf("не удалось получить токен: %s", err)
		log.Fatal(errMsg)
	}

	// campaignService := api_ozon.NewCampaign(apiService)
	campaignProductService := ozonAPI.NewCampaignProduct(apiService)

	products, err := campaignProductService.GetCampaignProductSku("3636696")
	products_, err := campaignProductService.Get("3636696")
	fmt.Printf("Количество продуктов в ответе products: %v\n\n", len(products))
	fmt.Printf("Товары в ответе: %v\n\n", products_)
	fmt.Printf("Товары в ответе skus: %v\n", products)
	fmt.Print("\n\n\n")

	products2, err := campaignProductService.GetCampaignProductSku("3636696")
	products2_, err := campaignProductService.Get("3636696")
	fmt.Printf("Количество продуктов в ответе products2: %v\n\n", len(products2))
	fmt.Printf("Товары в ответе: %v\n\n", products2_)
	fmt.Printf("Товары в ответе products2 skus: %v\n", products2)
	fmt.Print("\n\n\n")

	products3, err := campaignProductService.GetCampaignProductSku("3636696")
	products3_, err := campaignProductService.Get("3636696")
	fmt.Printf("Количество продуктов в ответе products: %v\n\n", len(products3))
	fmt.Printf("Товары в ответе: %v\n\n", products3_)
	fmt.Printf("Товары в ответе products3 skus: %v\n", products3)

	// Удалить все товары из рекламной компании.
	// err = RemoveAllProductFromCampaign(campaignProductService, "3636696")
	// fmt.Println(err)
}

func RemoveAllProductFromCampaign(service *ozonAPI.CampaignProduct, campaignId string) error {
	totalRemovedProducts := 0
	for {
		products, _ := service.GetCampaignProductSku(campaignId)
		if len(products) == 0 {
			break
		}
		if err := service.Remove(campaignId, products); err != nil {
			errMsg := fmt.Errorf("не удалось добавить товары в рекламную компанию: %s", err)
			log.Println(errMsg)
			continue
		}
		totalRemovedProducts += len(products)
	}
	log.Printf("Удалено товаров товаров - %v из рекламной компании: %s, клиент %s", totalRemovedProducts, campaignId, service.Client.ApiIdPerformance)
	return nil

}

func ChunkSlice(slice []string, chunkSize int) [][]string {
	var chunks [][]string
	for i := 0; i < len(slice); i += chunkSize {
		end := i + chunkSize
		if end > len(slice) {
			end = len(slice)
		}
		chunks = append(chunks, slice[i:end])
	}
	return chunks
}
