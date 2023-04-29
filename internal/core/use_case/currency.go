package useCase

import (
	"context"

	"github.com/isaias-dgr/currency-tracker/internal/core/domain"
	"github.com/isaias-dgr/currency-tracker/internal/core/ports"
	"go.uber.org/zap"
)

type Currency struct {
	repository ports.CurrencyRepository
	log        *zap.SugaredLogger
}

func NewCurrency(currency_repo ports.CurrencyRepository,
	log *zap.SugaredLogger) *Currency {
	return &Currency{
		repository: currency_repo,
		log:        log,
	}
}

func (c *Currency) GetByCode(ctx context.Context, code string, filter domain.Filter) (ta []*domain.CurrencyRepository, err error) {
	c.log.Infof("🐕‍🦺 Get Currencies by code %s", code)
	return c.repository.GetByCode(ctx, code, filter)
}
