package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/isaias-dgr/currency-tracker/src/currency/deliver/http"
	_CurrencyRepo "github.com/isaias-dgr/currency-tracker/src/currency/repository/postgres"
	useCase "github.com/isaias-dgr/currency-tracker/src/currency/usecase"
	"github.com/isaias-dgr/currency-tracker/src/domain"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

func SetUpLog() *zap.SugaredLogger {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	return logger.Sugar()
}

func SetUpRepository(logger *zap.SugaredLogger) (*sql.DB, domain.CurrencyRepository) {
	logger.Info("💾 Set up Database.")
	logger.Info(os.Getenv("MYSQL_USER"))
	connection := fmt.Sprintf(
		"postgresql://%s:%s@%s:%s?sslmode=disable",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"),
	)

	logger.Info(connection)
	dbConn, err := sql.Open(`postgres`, connection)
	if err != nil {
		logger.Error(err)
	}
	err = dbConn.Ping()
	if err != nil {
		logger.Error(err)
	}
	return dbConn, _CurrencyRepo.NewcurrencyRepository(dbConn, logger)
}

func main() {
	log := SetUpLog()
	msg := fmt.Sprintf(
		"🤓 SetUp %s_%s:%s..",
		os.Getenv("PROJ_NAME"),
		os.Getenv("PROJ_ENV"),
		"8080")
	log.Info(msg)
	log.Info("🚀 API V1.")
	dbConn, currency_repo := SetUpRepository(log)
	defer func() {
		err := dbConn.Close()
		if err != nil {
			log.Error(err)
		}
	}()
	useCase := useCase.NewCurrencyUseCase(currency_repo)
	http.NewCurrencyHandler(useCase, log)
}
