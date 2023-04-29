package useCase

import (
	"context"

	"github.com/isaias-dgr/currency-tracker/internal/core/ports"
)

type CurrencyIngest struct {
	repository ports.CurrencyRepository
	service    ports.CurrencyAPIService
}

func NewCurrencyIngest(
	currency_repo ports.CurrencyRepository,
	currency_service ports.CurrencyAPIService) *CurrencyIngest {
	return &CurrencyIngest{
		repository: currency_repo,
		service:    currency_service,
	}
}

func (c *CurrencyIngest) Run(ctx context.Context) (err error) {
	currencies_values, err := c.service.Get(ctx)
	if err != nil {
		return err
	}
	err = c.repository.InsertBulk(ctx, currencies_values)
	if err != nil {
		return err
	}
	return nil
}
