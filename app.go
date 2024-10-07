package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"wb-supply/internal/storage/sqlite"
	"wb-supply/internal/transformData"
	wbApi "wb-supply/internal/wbApi"
)

// App struct
type App struct {
	ctx      context.Context
	storage  *sqlite.Storage
	wbClient *wbApi.Client
}

type ExportWbData struct {
	Articles string `json:"articles"`
	Stats    string `json:"stats"`
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	storage, err := sqlite.New("./database/storage.db")
	if err != nil {
		panic(err)
	}
	a.storage = storage
}

func (a *App) CreateApiClient(apiKey string) error {
	wbClient, err := wbApi.New(apiKey)
	if err != nil {
		return err
	}
	a.wbClient = wbClient
	return nil
}

func (a *App) DeleteApiClient() {
	a.wbClient = nil
}

func (a *App) SaveApiKey(apiKey, apiKeyName string) (int64, error) {
	keyID, err := a.storage.AddApiKey(apiKey, apiKeyName)
	if err != nil {
		return 0, fmt.Errorf("error occured while saving key to db %v", err)
	}
	return keyID, nil
}

func (a *App) GetAllApiKeys() (apiKeys []*sqlite.ApiKey, e error) {
	apiKeys, e = a.storage.GetApiKeys()
	return
}

func (a *App) DeleteApiKey(apiKeyID string) error {
	keyId, err := strconv.ParseInt(apiKeyID, 10, 64)
	if err != nil {
		return fmt.Errorf("error occured while converting apiKeyID to int %v", err)
	}
	err = a.storage.DeleteApiKey(keyId)
	if err != nil {
		return fmt.Errorf("error occured while deleting api key %v", err)
	}
	return nil
}

func (a *App) GetWildberriesData(dateFrom string) (*ExportWbData, error) {
	if err := a.checkApiClient(); err != nil {
		return nil, err
	}
	stocks, err := a.wbClient.GetStocks(dateFrom)
	if err != nil {
		return nil, fmt.Errorf("error occured while getting wildberries data %v", err)
	}
	sales, err := a.wbClient.GetSales(dateFrom)
	if err != nil {
		return nil, fmt.Errorf("error occured while getting wildberries data %v", err)
	}
	if sales == nil {
		return nil, fmt.Errorf("no sales data %v", sales)
	}
	if stocks == nil {
		return nil, fmt.Errorf("no stocks data %v", stocks)
	}
	dataToTransform := &transformData.WbData{
		Sales:  &sales,
		Stocks: &stocks,
	}

	result, err := dataToTransform.GetStats()
	if err != nil {
		return nil, fmt.Errorf("error occured while creating stats %v", err)
	}

	exportArticles, err := json.Marshal(result.Articles)
	if err != nil {
		return nil, fmt.Errorf("error occured while creating export articles %v", err)
	}
	exportStats, err := json.Marshal(result.Stats)
	if err != nil {
		return nil, fmt.Errorf("error occured while creating export stats %v", err)
	}

	export := &ExportWbData{
		Articles: string(exportArticles),
		Stats:    string(exportStats),
	}

	return export, nil
}

func (a *App) checkApiClient() error {
	if a.wbClient == nil {
		return fmt.Errorf("api client hasn't created %v", a.wbClient)
	}
	return nil
}
