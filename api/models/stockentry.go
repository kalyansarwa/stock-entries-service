package models

import (
	"errors"
	"html"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

type StockEntry struct {
	ID                   uuid.UUID `json:"id"`
	PortfolioId          string    `gorm:"size:25;primary_key" json:"portfolioId"`
	Symbol               string    `gorm:"size:25;primary_key" json:"symbol"`
	EntryDate            time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"entryDate"`
	PurchasePricePerUnit float64   `gorm:"check:purchasePricePerUnit>0" json:"purchasePricePerUnit"`
	Quantity             uint      `gorm:"check:quantity > 0" json:"quantity"`
	PurchaseBrokerage    float64   `gorm:"check:purchaseBrokerage>0" json:"purchaseBrokerage"`
	TotalPurchasePrice   float64   `gorm:"check:totalPurchasePrice>0" json:"totalPurchasePrice"`
	Holding              bool      `gorm:"default:true" json:"holding"`
}

// Provides the default tablename for the structure
func (se *StockEntry) TableName() string {
	return "stocks.entries"
}

// Triggered before DB operation to create an entry
func (stockEntry *StockEntry) BeforeCreate(scope *gorm.Scope) error {
	scope.SetColumn("ID", uuid.New())
	return nil
}

// Similar to a constructor
func (se *StockEntry) Prepare() {
	se.PortfolioId = html.EscapeString(strings.TrimSpace(se.PortfolioId))
	se.Symbol = html.EscapeString(strings.TrimSpace(se.Symbol))
}

// Validation function - can be called to validate before any operation
func (se *StockEntry) Validate(action string) error {
	switch strings.ToLower(action) {
	case "update":
		if se.EntryDate.After(time.Now()) {
			return errors.New("entry date cannot be in future")
		}
		return nil
	default:
		return nil
	}
}

func (se *StockEntry) FindAllStockEntries(db *gorm.DB) (*[]StockEntry, error) {
	// var err error
	stockEntries := []StockEntry{}
	var err = db.Debug().Model(&StockEntry{}).Limit(100).Find(&stockEntries).Error
	if err != nil {
		return &[]StockEntry{}, err
	}
	return &stockEntries, nil
}

func (se *StockEntry) FindAllStockEntriesByPortfolioId(db *gorm.DB, portfolioId string) (*[]StockEntry, error) {
	// var err error
	stockEntries := []StockEntry{}
	var err = db.Debug().Model(&StockEntry{}).Where("portfolio_id = ?", portfolioId).Take(&stockEntries).Error
	if err != nil {
		return &[]StockEntry{}, err
	}
	if gorm.IsRecordNotFoundError(err) {
		return &[]StockEntry{}, errors.New("portfolio id invalid")
	}
	return &stockEntries, nil
}

func (se *StockEntry) FindStockEntryBySymbol(db *gorm.DB, portfolioId string, symbol string) (*StockEntry, error) {
	// var err error
	var err = db.Debug().Model(StockEntry{}).Where("portfolio_id = ? and symbol = ?", portfolioId, symbol).Take(&se).Error
	if err != nil {
		return &StockEntry{}, err
	}
	if gorm.IsRecordNotFoundError(err) {
		return &StockEntry{}, errors.New("stock symbol not found")
	}
	return se, err
}

func (se *StockEntry) SaveStockEntry(db *gorm.DB) (*StockEntry, error) {
	// var err error
	var err = db.Debug().Create(&se).Error
	if err != nil {
		return &StockEntry{}, err
	}
	return se, nil
}

func (se *StockEntry) UpdateStockEntry(db *gorm.DB, portfolioId string, symbol string) (*StockEntry, error) {

	db = db.Debug().Model(&StockEntry{}).Where("portfolio_id = ? and symbol = ?", portfolioId, symbol).Take(&StockEntry{}).UpdateColumns(
		map[string]interface{}{
			"entry_date":              se.EntryDate,
			"purchase_price_per_unit": se.PurchasePricePerUnit,
			"quantity":                se.Quantity,
			"purchase_brokerage":      se.PurchaseBrokerage,
			"total_purchase_price":    se.TotalPurchasePrice,
		},
	)
	if db.Error != nil {
		return &StockEntry{}, db.Error
	}

	// This is to display the updated user
	var err = db.Debug().Model(&StockEntry{}).Where("portfolio_id = ? and symbol = ?", portfolioId, symbol).Take(&se).Error
	if err != nil {
		return &StockEntry{}, err
	}
	return se, nil
}

func (se *StockEntry) DeleteStockEntry(db *gorm.DB, portfolioId string, symbol string) (int64, error) {
	db = db.Debug().Model(&StockEntry{}).Where("portfolio_id = ? and symbol = ?", portfolioId, symbol).Take(&StockEntry{}).Delete(&StockEntry{})
	if db.Error != nil {
		return 0, db.Error
	}
	return db.RowsAffected, nil
}

func (se *StockEntry) CountStockEntries(db *gorm.DB, portfolioId string) (int64, error) {

	var count int64;
	db = db.Debug().Model(&StockEntry{}).Where("portfolio_id = ?", portfolioId).Count(&count)
	if db.Error != nil {
		return 0, db.Error
	}
	return count, nil
}
