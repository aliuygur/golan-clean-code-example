package currencyratesapi

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/alioygur/golang-clean-code-example/app"
	"github.com/alioygur/golang-clean-code-example/currencyrates"
	"github.com/gorilla/mux"
)

type currencyratesService interface {
	FetchRates(ctx context.Context, params *currencyrates.FetchParams) (*currencyrates.RatesResponse, error)
	Exchange(ctx context.Context, amount float64, from currencyrates.CurrencyCode, to currencyrates.CurrencyCode) (float64, error)
}

type Handler struct {
	CurrencyRates currencyratesService
}

func (h *Handler) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/rates", h.getAllCurrencyRates)
	r.HandleFunc("/rates/exchange/{amount}/{from}/{to}", h.exchangeCurrency)
}

func (h *Handler) getAllCurrencyRates(w http.ResponseWriter, r *http.Request) {
	res, err := h.CurrencyRates.FetchRates(r.Context(), &currencyrates.FetchParams{Base: currencyrates.TRY})
	if err != nil {
		if errors.Is(err, currencyrates.ErrInvalidCurrencyCode) {
			app.BadRequest(w, r, err)
			return
		}
		app.InternalError(w, r, err)
		return
	}

	app.JSON(w, http.StatusOK, newRatesResponse(res))
}

func (h *Handler) exchangeCurrency(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	from, to := params["from"], params["to"]

	from, to = strings.ToUpper(from), strings.ToUpper(to)

	amount, err := strconv.ParseFloat(params["amount"], 64)
	if err != nil {
		app.BadRequest(w, r, fmt.Errorf("invalid amount: %s", params["amount"]))
		return
	}

	total, err := h.CurrencyRates.Exchange(r.Context(), amount, currencyrates.CurrencyCode(from), currencyrates.CurrencyCode(to))
	if err != nil {
		if errors.Is(err, currencyrates.ErrInvalidCurrencyCode) {
			app.BadRequest(w, r, err)
			return
		}
		app.InternalError(w, r, err)
		return
	}

	app.JSON(w, http.StatusOK, newExchangeResponse(total))
}
