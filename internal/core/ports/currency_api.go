package ports

import (
	"context"

	"github.com/isaias-dgr/currency-tracker/internal/core/domain"
)

type CurrencyAPIService interface {
	Get(context.Context) (domain.Currencies, error)
}
