package seed

import (
	"fmt"
	"os"
	"strconv"
	"time"
	"log"

	"github.com/gocarina/gocsv"
	"github.com/jinzhu/gorm"
	"github.com/kalyansarwa/stocksapi/api/models"
)

type FileEntry struct {
	EntryDate            string `csv:"entry_date"`
	Symbol               string `csv:"symbol"`
	PurchasePricePerUnit string `csv:"purchase_price_per_unit"`
	Quantity             string `csv:"quantity"`
	PurchaseBrokerage    string `csv:"purchase_brokerage"`
	TotalPurchasePrice   string `csv:"total_purchase_price"`
}

func LoadDB(db *gorm.DB) {

	stockEntry := models.StockEntry{}

	count, err := stockEntry.CountStockEntries(db, "skalyan")
	log.Printf("Database contains %d entries; Error: %s", count, err)

	if err == nil && count == 0 {

		entriesFile, err := os.OpenFile("stockEntries.csv", os.O_RDWR|os.O_CREATE, os.ModePerm)
		if err != nil {
			panic(err)
		}
		defer entriesFile.Close()

		stockEntries := []*FileEntry{}

		if err := gocsv.UnmarshalFile(entriesFile, &stockEntries); err != nil {
			panic(err)
		}

		for _, entry := range stockEntries {
			fmt.Println("Entry : ", entry)
			stockEntry := convertToStockEntry(entry)
			stockEntry.SaveStockEntry(db)
		}
	} 
}

func convertToStockEntry(entry *FileEntry) *models.StockEntry {

	stockEntry := models.StockEntry{}

	if entryDate, err := time.Parse("2-Jan-2006", entry.EntryDate); err == nil {
		stockEntry.EntryDate = entryDate
	}

	stockEntry.Symbol = entry.Symbol

	if s, err := strconv.ParseFloat(entry.PurchasePricePerUnit, 64); err == nil {
		stockEntry.PurchasePricePerUnit = s
	}

	if s, err := strconv.ParseUint(entry.Quantity, 10, 32); err == nil {
		stockEntry.Quantity = uint(s)
	}

	if s, err := strconv.ParseFloat(entry.PurchaseBrokerage, 64); err == nil {
		stockEntry.PurchaseBrokerage = s
	}

	if s, err := strconv.ParseFloat(entry.TotalPurchasePrice, 64); err == nil {
		stockEntry.TotalPurchasePrice = s
	}

	stockEntry.PortfolioId = "skalyan"

	return &stockEntry

}
