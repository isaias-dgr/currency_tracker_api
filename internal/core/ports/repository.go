package ports

import (
	"context"

	"github.com/isaias-dgr/currency-tracker/internal/core/domain"
)

type CurrencyRepository interface {
	InsertBulk(ctx context.Context, ta domain.Currencies) error
	GetByCode(context.Context, string, domain.Filter) ([]*domain.CurrencyRepository, error)
}
