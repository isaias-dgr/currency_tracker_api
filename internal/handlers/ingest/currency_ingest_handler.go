package ingest

import (
	"context"
	"time"

	useCase "github.com/isaias-dgr/currency-tracker/internal/core/use_case"
	"github.com/newrelic/go-agent/v3/integrations/nrzap"
	"github.com/newrelic/go-agent/v3/newrelic"
	"go.uber.org/zap"
)

type CurrencyIngestHandler struct {
	Ingest    *useCase.CurrencyIngest
	log       *zap.SugaredLogger
	sleeptime int
}

func NewCurrencyIngestHandler(CurrencyIngest *useCase.CurrencyIngest,
	logger *zap.SugaredLogger,
	sleeptime int) {
	handler := &CurrencyIngestHandler{
		Ingest:    CurrencyIngest,
		log:       logger,
		sleeptime: sleeptime,
	}
	handler.log.Infoln("Ingest 🤮")
	app, _ := newrelic.NewApplication(
		newrelic.ConfigAppName("currency-ingest"),
		newrelic.ConfigLicense("918bbc0a76f525bfbdbbba48f6c7b2831019NRAL"),
		newrelic.ConfigAppLogForwardingEnabled(true),
		nrzap.ConfigLogger(handler.log.Desugar().Named("newrelic")),
	)

	ctx := context.Background()
	for range time.Tick(time.Minute * time.Duration(sleeptime)) {
		func() {
			txn := app.StartTransaction("currency-ingest_request")
			defer txn.End()
			if err := handler.Ingest.Run(ctx); err != nil {
				handler.log.Errorln("Ingest 🤮")
			}
		}()

	}

	handler.log.Infoln("End Ingest 🤮")
}
