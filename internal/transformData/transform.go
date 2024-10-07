package transformData

import (
	"maps"
	"sync"
	"wb-supply/internal/wbApi"
)

//go:generate go run github.com/vektra/mockery/v2@v2.46.2

type Warehouse[wh string, count int] map[wh]count
type Counted[art string, wh Warehouse[string, int]] map[art]wh

type Articles map[string]map[string]bool
type Stats map[string]map[string]map[string]int

type WbData struct {
	Articles      Articles       `json:"articles"`
	Stocks        *[]wbApi.Stock `json:"stocks"`
	Sales         *[]wbApi.Sale  `json:"sales"`
	CountedSales  *Counted[string, Warehouse[string, int]]
	CountedStocks *Counted[string, Warehouse[string, int]]
}

type Export struct {
	Articles Articles `json:"articles"`
	Stats    Stats    `json:"stats"`
}

func (a Articles) addArticle(article string, warehouse string) {
	if a[article] == nil {
		a[article] = make(map[string]bool)
	}
	a[article][warehouse] = true
}

func renameWarehouse(warehouse string) string {
	switch warehouse {
	case
		"Электросталь",
		"Коледино",
		"Тула":
		return "Москва"
	case
		"Санкт-Петербург Уткина Заводь":
		return "СПб"
	default:
		return warehouse
	}
}

func (d *WbData) GetStocksCountByWarehouse(mutex *sync.Mutex) {
	countedStocks := make(Counted[string, Warehouse[string, int]])
	stocks := d.Stocks
	articles := make(Articles)

	for _, stock := range *stocks {
		supplierArticle := *stock.SupplierArticle
		warehouse := *stock.WarehouseName

		quantity := *stock.Quantity
		if countedStocks[supplierArticle] == nil {
			countedStocks[supplierArticle] = make(Warehouse[string, int])
		}
		warehouse = renameWarehouse(warehouse)
		countedStocks[supplierArticle][warehouse] += quantity
		articles.addArticle(supplierArticle, warehouse)
	}
	d.copyArticlesToWbData(mutex, articles)
	d.CountedStocks = &countedStocks
}

func (d *WbData) GetSalesCountByWarehouse(mutex *sync.Mutex) {
	countedSales := make(Counted[string, Warehouse[string, int]])
	sales := d.Sales
	articles := make(Articles)

	for _, sale := range *sales {
		supplierArticle := *sale.SupplierArticle
		warehouse := *sale.WarehouseName

		if countedSales[supplierArticle] == nil {
			countedSales[supplierArticle] = make(Warehouse[string, int])
		}
		warehouse = renameWarehouse(warehouse)
		articles.addArticle(supplierArticle, warehouse)
		countedSales[supplierArticle][warehouse] += 1
	}
	d.copyArticlesToWbData(mutex, articles)
	d.CountedSales = &countedSales
}

func (d *WbData) copyArticlesToWbData(mutex *sync.Mutex, articles Articles) {
	mutex.Lock()
	if d.Articles == nil {
		d.Articles = make(Articles)
	}
	maps.Copy(d.Articles, articles)
	mutex.Unlock()
}

func (d *WbData) GetStats() (*Export, error) {
	var wg sync.WaitGroup
	var mutex sync.Mutex
	wg.Add(2)
	go func() {
		defer wg.Done()
		d.GetStocksCountByWarehouse(&mutex)
	}()
	go func() {
		defer wg.Done()
		d.GetSalesCountByWarehouse(&mutex)
	}()
	wg.Wait()

	articles := d.Articles
	result := make(Stats)

	for article, warehouses := range articles {
		result[article] = make(map[string]map[string]int)
		for warehouse := range warehouses {
			result[article][warehouse] = make(map[string]int)
			sales := (*d.CountedSales)[article][warehouse]
			stocks := (*d.CountedStocks)[article][warehouse]
			result[article][warehouse]["sales"] = sales
			result[article][warehouse]["stocks"] = stocks
			result[article][warehouse]["supply"] = sales - stocks
		}
	}

	return &Export{
		Articles: articles,
		Stats:    result,
	}, nil
}
