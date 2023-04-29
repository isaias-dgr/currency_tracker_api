package repositories

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/isaias-dgr/currency-tracker/internal/core/domain"
	"github.com/isaias-dgr/currency-tracker/internal/core/ports"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
)

type SuiteRepository struct {
	suite.Suite
	db      *sql.DB
	mockSQL sqlmock.Sqlmock
	repo    ports.CurrencyRepository
}

func (s *SuiteRepository) SetupTest() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	sugar := logger.Sugar()

	db, mockSQL, err := sqlmock.New()
	if err != nil {
		s.Failf("an error '%s' was not expected when opening a stub database connection", err.Error())
	}
	s.db = db
	s.mockSQL = mockSQL
	s.repo = NewCurrencyRepo(db, sugar)
}

func (s *SuiteRepository) TestGetByCode() {
	s.Run("Success test return a Currencies by Code", func() {

		raw_uuid := uuid.New()
		binary_uuid, _ := raw_uuid.MarshalBinary()
		code := "MXN"
		rows := []string{"id", "code", "value", "updated_at", "created_at"}
		created_at := time.Now()
		mockTask := domain.Currency{
			Code:  code,
			Value: 19.342,
		}
		filter := domain.Filter{
			Offset: 0,
			Limit:  10,
		}
		data := sqlmock.NewRows(rows).
			AddRow(binary_uuid, created_at, created_at,
				mockTask.Code, mockTask.Value)

		q := "SELECT \\* FROM currency WHERE code=\\$1 ORDER BY created_at ASC LIMIT \\$2 OFFSET \\$3"
		s.mockSQL.ExpectQuery(q).
			WithArgs(code, filter.Limit, filter.Offset).
			WillReturnRows(data)

		task, err := s.repo.GetByCode(context.TODO(), code, filter)

		s.NoError(err)
		s.Equal(mockTask.Code, task[0].Code)
		s.Equal(mockTask.Value, task[0].Value)
	})

	s.Run("Success test return a Currencies by Code and range dates", func() {
		raw_uuid := uuid.New()
		binary_uuid, _ := raw_uuid.MarshalBinary()
		code := "DLL"
		rows := []string{"id", "code", "value", "updated_at", "created_at"}
		created_at := time.Now()
		mockTask := domain.Currency{
			Code:  code,
			Value: 19.342,
		}
		filter := domain.Filter{
			Offset: 0,
			Limit:  10,
			Finit:  &created_at,
			Fend:   &created_at,
		}
		data := sqlmock.NewRows(rows).
			AddRow(binary_uuid, created_at, created_at,
				mockTask.Code, mockTask.Value)

		q := "SELECT \\* FROM currency WHERE code=\\$1 AND created_at > \\$2 AND created_at < \\$3 ORDER BY created_at ASC LIMIT \\$4 OFFSET \\$5"
		s.mockSQL.ExpectQuery(q).
			WithArgs(code, filter.Finit, filter.Fend, filter.Limit, filter.Offset).
			WillReturnRows(data)

		task, err := s.repo.GetByCode(context.TODO(), code, filter)

		s.NoError(err)
		s.Equal(mockTask.Code, task[0].Code)
		s.Equal(mockTask.Value, task[0].Value)
	})

	s.Run("When the query not found task", func() {
		code := "BTC"
		rows := []string{"id", "code", "value", "updated_at", "created_at"}
		filter := domain.Filter{
			Offset: 0,
			Limit:  10,
		}
		data := sqlmock.NewRows(rows)

		q := "SELECT \\* FROM currency WHERE code=\\$1 ORDER BY created_at ASC LIMIT \\$2 OFFSET \\$3"
		s.mockSQL.ExpectQuery(q).
			WithArgs(code, filter.Limit, filter.Offset).
			WillReturnRows(data)

		task, err := s.repo.GetByCode(context.TODO(), code, filter)

		s.Error(err)
		s.Equal("not_found", err.Error())
		s.Nil(task)
	})

	s.Run("When test code without format return error", func() {
		task, err := s.repo.GetByCode(context.TODO(), "", domain.Filter{})
		s.Error(err)
		s.Nil(task)
	})

	s.Run("When the query fails return error", func() {
		q := "SELECT \\* FROM currency WHERE code=\\$1 ORDER BY created_at ASC LIMIT \\$2 OFFSET \\$3"
		filter := domain.Filter{
			Offset: 0,
			Limit:  10,
		}
		s.mockSQL.ExpectQuery(q).
			WithArgs("MXN", filter.Limit, filter.Offset).
			WillReturnError(errors.New("generic error"))

		task, err := s.repo.GetByCode(context.TODO(), "", domain.Filter{})
		s.Error(err)
		s.Equal("query_context", err.Error())
		s.Nil(task)
	})
}

func (s *SuiteRepository) TestInsert() {
	s.Run("Success test return a cyrrency", func() {
		currency := domain.NewCurrency("ETH", 200.3)
		q := "INSERT INTO currency\\(code, value, created_at, updated_at\\) VALUES \\(\\$1,\\$2,\\$3,\\$4\\);"
		s.mockSQL.
			ExpectPrepare(q).
			ExpectExec().
			WithArgs(currency.Code, currency.Value,
				sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := s.repo.InsertBulk(context.TODO(), domain.Currencies{*currency})
		s.Nil(err)
	})

	s.Run("Success test return a cyrrency", func() {
		currency_1 := domain.NewCurrency("COL", 200.3)
		currency_2 := domain.NewCurrency("ARG", 200.3)
		q := "INSERT INTO currency\\(code, value, created_at, updated_at\\) VALUES \\(\\$1,\\$2,\\$3,\\$4\\),\\(\\$5,\\$6,\\$7,\\$8\\);"
		s.mockSQL.
			ExpectPrepare(q).
			ExpectExec().
			WithArgs(
				currency_1.Code, currency_2.Value, sqlmock.AnyArg(), sqlmock.AnyArg(),
				currency_2.Code, currency_2.Value, sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := s.repo.InsertBulk(context.TODO(),
			domain.Currencies{*currency_1, *currency_2})
		s.Nil(err)
	})

	s.Run("When the prepare context faild must return error", func() {
		currency := domain.NewCurrency("ETH", 200.3)
		q := "INSERT INTO currency\\(code, value, created_at, updated_at\\) VALUES \\(\\$1,\\$2,\\$3,\\$4\\);"
		s.mockSQL.
			ExpectPrepare(q).
			WillReturnError(errors.New("prepare error"))

		err := s.repo.InsertBulk(context.TODO(), domain.Currencies{*currency})
		s.Error(err)
		s.Equal("query_prepare_ctx", err.Error())
	})

	s.Run("When the Exec stmt faild must return error", func() {
		currency := domain.NewCurrency("ETH", 200.3)
		q := "INSERT INTO currency\\(code, value, created_at, updated_at\\) VALUES \\(\\$1,\\$2,\\$3,\\$4\\);"
		s.mockSQL.
			ExpectPrepare(q).
			ExpectExec().
			WithArgs(currency.Code, currency.Value,
				sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnError(errors.New("exec error"))

		err := s.repo.InsertBulk(context.TODO(), domain.Currencies{*currency})

		s.Error(err)
		s.Equal("query_exec", err.Error())
	})

	s.Run("When the Exec result send error must return error", func() {
		currency := domain.NewCurrency("ETH", 200.3)
		q := "INSERT INTO currency\\(code, value, created_at, updated_at\\) VALUES \\(\\$1,\\$2,\\$3,\\$4\\);"
		s.mockSQL.ExpectPrepare(q).
			ExpectExec().
			WithArgs(currency.Code, currency.Value,
				sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnResult(sqlmock.NewErrorResult(errors.New("not_found")))

		err := s.repo.InsertBulk(context.TODO(), domain.Currencies{*currency})
		s.NotNil(err)
		s.Equal("query_exec", err.Error())
	})
}

func TestSuiteRepository(t *testing.T) {
	suite.Run(t, new(SuiteRepository))
}
