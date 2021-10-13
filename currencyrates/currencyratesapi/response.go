package currencyratesapi

import (
	"github.com/alioygur/golang-clean-code-example/currencyrates"
)

type ratesResponse struct {
	Base  string             `json:"base"`
	Rates map[string]float64 `json:"rates"`
}

func newRatesResponse(res *currencyrates.RatesResponse) *ratesResponse {
	ratesRes := ratesResponse{
		Base:  string(res.Base),
		Rates: make(map[string]float64),
	}

	for code, rate := range res.Rates {
		ratesRes.Rates[string(code)] = rate
	}
	return &ratesRes
}

type exchangeResponse struct {
	Result float64 `json:"result"`
}

func newExchangeResponse(result float64) *exchangeResponse {
	return &exchangeResponse{Result: result}
}
