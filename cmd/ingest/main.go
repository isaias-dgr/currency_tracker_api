package main

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"

	useCase "github.com/isaias-dgr/currency-tracker/internal/core/use_case"
	"github.com/isaias-dgr/currency-tracker/internal/external_api"
	"github.com/isaias-dgr/currency-tracker/internal/handlers/ingest"
	"github.com/isaias-dgr/currency-tracker/internal/repositories"

	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

func SetUpLog() *zap.SugaredLogger {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	return logger.Sugar()
}

func SetUpRepository(logger *zap.SugaredLogger) (*sql.DB, *repositories.Currencies_psq) {
	logger.Info("💾 Set up Database.")
	logger.Info(os.Getenv("POSTGRES_USER"))
	connection := fmt.Sprintf(
		"postgresql://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_DATABASE"),
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
	return dbConn, repositories.NewCurrencyRepo(dbConn, logger)
}

func SetUPCurrencyAPI(logger *zap.SugaredLogger) *external_api.CurrencyAPI {
	logger.Info("🤑 Set up currency api")
	time_out, _ := strconv.Atoi(os.Getenv("TIMEOUT"))
	return external_api.NewCurrencyAPI(
		os.Getenv("CURRENCY_URL"),
		os.Getenv("CURRENCY_APIKEY"),
		time_out,
	)
}

func main() {
	log := SetUpLog()
	log.Info("🤓 Ingest")
	dbConn, currency_repo := SetUpRepository(log)
	defer func() {
		err := dbConn.Close()
		if err != nil {
			log.Error(err)
		}
	}()

	currency_api := SetUPCurrencyAPI(log)
	use_case := useCase.NewCurrencyIngest(currency_repo, currency_api)
	sleep_time, err := strconv.Atoi(os.Getenv("SLEEP_TIME"))
	if err != nil {
		log.Error(err)
	}
	ingest.NewCurrencyIngestHandler(use_case, log, sleep_time)
}
