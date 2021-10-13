package currencyratesapi

import (
	"context"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/alioygur/golang-clean-code-example/currencyrates"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockService struct {
	mock.Mock
}

func (m *MockService) FetchRates(ctx context.Context, params *currencyrates.FetchParams) (*currencyrates.RatesResponse, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*currencyrates.RatesResponse), args.Error(1)
}
func (m *MockService) Exchange(ctx context.Context, amount float64, from currencyrates.CurrencyCode, to currencyrates.CurrencyCode) (float64, error) {
	args := m.Called(ctx, amount, from, to)
	return args.Get(0).(float64), args.Error(1)
}

func Test_currencyratesHandler_exchangeCurrency(t *testing.T) {
	mockService := &MockService{}

	r := mux.NewRouter()

	handler := Handler{CurrencyRates: mockService}
	handler.RegisterRoutes(r)

	t.Run("success", func(t *testing.T) {
		mockService.
			On("Exchange", mock.AnythingOfType("*context.valueCtx"), 2.0, currencyrates.TRY, currencyrates.USD).
			Return(16.0, nil).
			Once()

		w := httptest.NewRecorder()

		r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/rates/exchange/2/TRY/USD", nil))

		res := w.Result()
		defer res.Body.Close()

		body, err := ioutil.ReadAll(res.Body)
		require.NoError(t, err)

		assert.Equal(t, res.StatusCode, http.StatusOK)
		assert.Equal(t, "application/json; charset=utf-8", res.Header.Get("content-type"))
		assert.Equal(t, `{"result":16}`, strings.TrimSuffix(string(body), "\n"))

		mockService.AssertExpectations(t)
	})

	t.Run("invalid currency code", func(t *testing.T) {
		mockService.
			On("Exchange", mock.AnythingOfType("*context.valueCtx"), 2.0, currencyrates.CurrencyCode("X"), currencyrates.USD).
			Return(0.0, currencyrates.ErrInvalidCurrencyCode).
			Once()

		w := httptest.NewRecorder()

		r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/rates/exchange/2/X/USD", nil))

		res := w.Result()
		defer res.Body.Close()

		body, err := ioutil.ReadAll(res.Body)
		require.NoError(t, err)

		assert.Equal(t, res.StatusCode, http.StatusBadRequest)
		assert.Equal(t, "application/json; charset=utf-8", res.Header.Get("content-type"))
		assert.Equal(t, `{"status":400,"code":"INVALID_CURRENCY_CODE","error":"invalid currency code"}`, strings.TrimSuffix(string(body), "\n"))

		mockService.AssertExpectations(t)
	})
}

func TestHandler_getAllCurrencyRates(t *testing.T) {
	mockService := &MockService{}

	r := mux.NewRouter()

	handler := Handler{CurrencyRates: mockService}
	handler.RegisterRoutes(r)

	mockService.
		On("FetchRates", mock.AnythingOfType("*context.valueCtx"), &currencyrates.FetchParams{Base: currencyrates.TRY}).
		Return(&currencyrates.RatesResponse{
			Base:  currencyrates.TRY,
			Rates: map[currencyrates.CurrencyCode]float64{currencyrates.TRY: 1, currencyrates.USD: 8},
		}, nil).
		Once()

	w := httptest.NewRecorder()

	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/rates", nil))

	res := w.Result()
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	require.NoError(t, err)

	assert.Equal(t, res.StatusCode, http.StatusOK)
	assert.Equal(t, "application/json; charset=utf-8", res.Header.Get("content-type"))
	assert.Equal(t, `{"base":"TRY","rates":{"TRY":1,"USD":8}}`, strings.TrimSuffix(string(body), "\n"))

	mockService.AssertExpectations(t)
}
