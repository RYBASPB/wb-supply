package wbApi

import (
	"context"
	"fmt"

	"github.com/oapi-codegen/oapi-codegen/v2/pkg/securityprovider"
)

const WB_STATS = "https://statistics-api.wildberries.ru/"

type Stock struct {
	Quantity        *int    `json:"quantity,omitempty"`
	SupplierArticle *string `json:"supplierArticle,omitempty"`
	WarehouseName   *string `json:"warehouseName,omitempty"`
}

type Sale struct {
	SupplierArticle *string `json:"supplierArticle,omitempty"`
	WarehouseName   *string `json:"warehouseName,omitempty"`
}

func New(apiKey string) (client *Client, err error) {
	auth, err := securityprovider.NewSecurityProviderApiKey("header", "Authorization", apiKey)
	if err != nil {
		return nil, err
	}
	client, err = NewClient(WB_STATS, WithRequestEditorFn(auth.Intercept))
	if err != nil {
		return nil, fmt.Errorf("falied to create client: %v", err)
	}
	return client, nil
}

func (c *Client) GetStocks(dateFrom string) (stocks []Stock, err error) {
	res, err := c.GetApiV1SupplierStocks(context.Background(), &GetApiV1SupplierStocksParams{
		DateFrom: dateFrom,
	})
	if err != nil {
		return nil, fmt.Errorf("falied to get stocks: %v", err)
	}
	parsedResponse, err := ParseGetApiV1SupplierStocksResponse(res)
	if err != nil {
		return nil, fmt.Errorf("falied to parse stocks: %v", err)
	}
	stocksFromResponse := parsedResponse.JSON200
	for _, stock := range *stocksFromResponse {
		requiredStock := Stock{
			Quantity:        stock.Quantity,
			SupplierArticle: stock.SupplierArticle,
			WarehouseName:   stock.WarehouseName,
		}
		stocks = append(stocks, requiredStock)
	}
	return
}

func (c *Client) GetSales(dateFrom string) (sales []Sale, err error) {
	res, err := c.GetApiV1SupplierSales(context.Background(), &GetApiV1SupplierSalesParams{
		DateFrom: dateFrom,
	})
	if err != nil {
		return nil, fmt.Errorf("falied to get stocks: %v", err)
	}
	parsedResponse, err := ParseGetApiV1SupplierSalesResponse(res)
	if err != nil {
		return nil, fmt.Errorf("falied to parse stocks: %v", err)
	}
	salesFromResponse := parsedResponse.JSON200
	for _, sale := range *salesFromResponse {
		requiredSale := Sale{
			SupplierArticle: sale.SupplierArticle,
			WarehouseName:   sale.WarehouseName,
		}
		sales = append(sales, requiredSale)
	}
	return
}
