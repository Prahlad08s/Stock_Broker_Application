package repositories

import (
	dbModel "stock_broker_application/src/models"

	"gorm.io/gorm"
)

type GetStocksSymbolRepository interface {
	GetStocksSymbols() ([]string, error)
}

type UserDbRepository struct {
	db *gorm.DB
}

func NewGetStocksSymbolRepository(db *gorm.DB) GetStocksSymbolRepository {
	return &UserDbRepository{db: db}
}

func (user *UserDbRepository) GetStocksSymbols() ([]string, error) {
	var symbols []string
	if err := user.db.Model(&dbModel.Stocks{}).Pluck("symbol", &symbols).Error; err != nil {
		return nil, err
	}
	return symbols, nil
}
