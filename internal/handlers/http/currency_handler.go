package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/isaias-dgr/currency-tracker/internal/core/domain"
	useCase "github.com/isaias-dgr/currency-tracker/internal/core/use_case"
	"github.com/newrelic/go-agent/v3/integrations/nrzap"
	"github.com/newrelic/go-agent/v3/newrelic"
	"go.uber.org/zap"
)

type CurrencyHandler struct {
	use_case useCase.Currency
	log      *zap.SugaredLogger
}

func NewCurrencyHandler(currencyUseCase useCase.Currency, logger *zap.SugaredLogger) {
	handler := &CurrencyHandler{
		use_case: currencyUseCase,
		log:      logger,
	}

	app, _ := newrelic.NewApplication(
		newrelic.ConfigAppName("currency-tracker"),
		newrelic.ConfigLicense(""),
		newrelic.ConfigAppLogForwardingEnabled(true),
		nrzap.ConfigLogger(handler.log.Desugar().Named("newrelic")),
	)
	r := mux.NewRouter()

	r.HandleFunc(newrelic.WrapHandleFunc(app, "/currencies", handler.GetCurrency)).Methods("GET")
	r.HandleFunc(newrelic.WrapHandleFunc(app, "/currencies/", handler.GetCurrency)).Methods("GET")
	r.HandleFunc(newrelic.WrapHandleFunc(app, "/currencies/{currency_code}", handler.GetCurrency)).Methods("GET")
	log.Fatal(http.ListenAndServe(":8080", r))
}

func (t *CurrencyHandler) DecoderBody(b io.ReadCloser, ta *domain.Currency) error {
	var unmarshalErr *json.UnmarshalTypeError
	decoder := json.NewDecoder(b)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(ta)
	if err != nil {
		t.log.Error("Bad Request. %s", err.Error())
		if errors.As(err, &unmarshalErr) {
			return fmt.Errorf("bad request:. Wrong Type provided for field %s", unmarshalErr.Field)
		} else {
			return fmt.Errorf("bad request: %s", err.Error())
		}
	}
	return nil
}

func validate(t *domain.Currency) error {
	if strings.TrimSpace(t.Code) == "" {
		return errors.New("bad request: Requiered field code")
	}

	return nil
}

func (c *CurrencyHandler) GetCurrency(w http.ResponseWriter, r *http.Request) {
	c.log.Infow("🕸 Get by code handler", "url", r.URL, "method", r.Method)
	vars := mux.Vars(r)
	code := "ALL"
	if val, ok := vars["currency_code"]; ok {
		code = val
	}
	filter := NewFilter(r.URL.Query())
	currencies, err := c.use_case.GetByCode(r.Context(), code, *filter)
	if err != nil {
		errorResponse(w, http.StatusNotFound, err.Error())
		return
	}
	makeResponse(w, http.StatusOK, currencies, filter, len(currencies))
}

func makeResponse(w http.ResponseWriter,
	code int,nc makeResponse(w http.ResponseWriter,
	code int,
	body interface{},
	filter *domain.Filter,
	total int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	resp := NewResponse(body, total, filter, "")
	jsonResp, _ := json.Marshal(resp)
	w.Write(jsonResp)
}

func errorResponse(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	resp := make(map[string]string)
	resp["message"] = message
	jsonResp, _ := json.Marshal(resp)
	body interface{},
	filter *domain.Filter,
	total int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	resp := NewResponse(body, total, filter, "")
	jsonResp, _ := json.Marshal(resp)
	w.Write(jsonResp)
}

func errorResponse(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	resp := make(map[string]string)
	resp["message"] = message
	jsonResp, _ := json.Marshal(resp)
	w.Write(jsonResp)
}
