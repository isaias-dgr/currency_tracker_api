package main

import (
	"database/sql"
	"fmt"
	"os"

	useCase "github.com/isaias-dgr/currency-tracker/internal/core/use_case"
	"github.com/isaias-dgr/currency-tracker/internal/handlers/http"
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

func main() {
	log := SetUpLog()
	msg := fmt.Sprintf(
		"🤓 SetUp cmd %s_%s:%s..",
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

	u := useCase.NewCurrency(currency_repo, log)
	http.NewCurrencyHandler(*u, log)
}
