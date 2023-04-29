package useCase

import (
	"context"
	"errors"
	"testing"

	"github.com/isaias-dgr/currency-tracker/internal/core/domain"
	mocks "github.com/isaias-dgr/currency-tracker/internal/core/ports/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type UseCaseIngestSuite struct {
	suite.Suite
	repo *mocks.CurrencyRepository
	service *mocks.CurrencyAPIService
	uCase   *CurrencyIngest
}

func (s *UseCaseIngestSuite) SetupTest() {
	s.repo = new(mocks.CurrencyRepository)
	s.service = new(mocks.CurrencyAPIService)
	s.uCase = NewCurrencyIngest(s.repo, s.service)
}

func (s *UseCaseIngestSuite) TestGetByCode() {
 	currencies := domain.Currencies{
		domain.Currency{
			Code: "MXN", 
			Value: 12.312,
		},
		domain.Currency{
			Code: "MXN", 
			Value: 12.312,
		},
	}

	s.service.On("Get", mock.Anything).Return(currencies, nil)
	s.repo.On("InsertBulk", mock.Anything, currencies).Return(nil)
	ctx := context.Background()
	err := s.uCase.Run(ctx)
	assert.Nil(s.T(), err, "The get mock its not working")
}

func (s *UseCaseIngestSuite) TestGetByCodeErrorService() {
	s.service.On("Get", mock.Anything).Return(nil, errors.New("Error general"))
	ctx := context.Background()
	err := s.uCase.Run(ctx)
	assert.NotNil(s.T(), err, "The get mock its not working")
	assert.Equal(s.T(), "Error general", err.Error())
}


func (s *UseCaseIngestSuite) TestGetByCodeErrorRepository() {
 	currencies := domain.Currencies{
		domain.Currency{
			Code: "MXN", 
			Value: 12.312,
		},
		domain.Currency{
			Code: "MXN", 
			Value: 12.312,
		},
	}

	s.service.On("Get", mock.Anything).Return(currencies, nil)
	s.repo.On("InsertBulk", mock.Anything, currencies).
		   Return(errors.New("Repo error"))
	ctx := context.Background()
	err := s.uCase.Run(ctx)
	assert.NotNil(s.T(), err, "The get mock its not working")
	assert.Equal(s.T(), "Repo error", err.Error())
}

func TestUseCaseIngestSuite(t *testing.T) {
	suite.Run(t, new(UseCaseIngestSuite))
}
