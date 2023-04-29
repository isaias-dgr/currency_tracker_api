package useCase

import (
	"context"
	"testing"

	"github.com/isaias-dgr/currency-tracker/internal/core/domain"
	mocks "github.com/isaias-dgr/currency-tracker/internal/core/ports/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
)

type UseCaseSuite struct {
	suite.Suite
	repo  *mocks.CurrencyRepository
	uCase *Currency
}

func (s *UseCaseSuite) SetupTest() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	sugar := logger.Sugar()

	s.repo = new(mocks.CurrencyRepository)
	s.uCase = NewCurrency(s.repo, sugar)
}

func (s *UseCaseSuite) TestGetByCode() {
	filter := domain.Filter{
		Limit:  100,
		Offset: 0,
	}
	s.repo.On("GetByCode", mock.Anything, "MXN", filter).Return(nil, nil)
	ctx := context.Background()

	_, err := s.uCase.GetByCode(ctx, "MXN", filter)
	assert.Nil(s.T(), err, "The get mock its not working")
}

func TestUseCaseSuite(t *testing.T) {
	suite.Run(t, new(UseCaseSuite))
}
