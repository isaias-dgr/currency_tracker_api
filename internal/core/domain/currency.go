package domain

import (
	"time"

	"github.com/google/uuid"
)

type Currency struct {
	Code  string
	Value float64
}

type Currencies []Currency

type CurrencyRepository struct {
	ID uuid.UUID `json:"id"`
	Currency
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
}

type Filter struct {
	Offset int
	Limit  int
	Finit  *time.Time
	Fend   *time.Time
	SortBy string
}

func NewCurrency(code string, value float64) *Currency {
	return &Currency{
		Code:  code,
		Value: value,
	}
}

type CurrenciesResponse struct {
	Data  []*Currency
	Total int
}

// func NewCurrencies(currency []*Currency, total int) *Currencies {
// 	return &Currencies{
// 		Data:  ts,
// 		Total: total,
// 	}
// }
