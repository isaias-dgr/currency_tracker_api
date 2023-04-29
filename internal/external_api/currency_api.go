package external_api

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/isaias-dgr/currency-tracker/internal/core/domain"
)

type ValueCurrencyAPI struct {
	Code  string  `json:"code"`
	Value float64 `json:"value"`
}

type MetaCurrencyAPI struct {
	LastUpdateAt string `json:"last_updated_at"`
}

type CurrenciesAPI struct {
	Meta MetaCurrencyAPI
	Data map[string]ValueCurrencyAPI `json:"data"`
}

type CurrencyAPI struct {
	client *http.Client
	url    string
	apikey string
}

func NewCurrencyAPI(url, apikey string, timeout int) *CurrencyAPI {
	return &CurrencyAPI{
		url:    fmt.Sprintf("%s?apikey=%s", url, apikey),
		client: &http.Client{
			Timeout: time.Duration(timeout) * time.Second,
		},
		apikey: apikey,
	}
}

func (c *CurrencyAPI) Get(ctx context.Context) (domain.Currencies, error) {
	request, err := http.NewRequest("GET", c.url, nil)
	if err != nil {
		return nil, err
	}

	response, err := c.client.Do(request)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	currencyResponse := &CurrenciesAPI{}
	if err := json.Unmarshal(body, currencyResponse); err != nil {
		return nil, err
	}

	currencies_api := domain.Currencies{}
	for _, value := range currencyResponse.Data {
		currencies_api = append(currencies_api, domain.Currency{
			Code:  value.Code,
			Value: value.Value,
		})
	}

	return currencies_api, nil
}
