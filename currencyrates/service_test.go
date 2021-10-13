package currencyrates

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockProvider struct {
	mock.Mock
}

func (m *MockProvider) FetchRates(ctx context.Context, params *FetchParams) (*RatesResponse, error) {
	args := m.Called(ctx, params)

	return args.Get(0).(*RatesResponse), args.Error(1)
}

func (m *MockProvider) GetRate(ctx context.Context, from, to CurrencyCode) (float64, error) {
	args := m.Called(ctx, from, to)
	return args.Get(0).(float64), args.Error(1)
}

func TestService_Exchange(t *testing.T) {
	mockProvider := new(MockProvider)
	mockProvider.
		On("FetchRates", context.Background(), &FetchParams{Base: TRY}).
		Return(&RatesResponse{
			Base:  TRY,
			Rates: map[CurrencyCode]float64{TRY: 1, USD: 8},
		}, nil).
		Once()

	service := NewService(mockProvider)

	total, err := service.Exchange(context.Background(), 2, TRY, USD)
	assert.NoError(t, err)
	assert.Equal(t, 16.0, total)

	mockProvider.AssertExpectations(t)

}

func TestService_FetchRates(t *testing.T) {
	mockProvider := new(MockProvider)
	mockProvider.
		On("FetchRates", context.Background(), &FetchParams{Base: TRY}).
		Return(&RatesResponse{
			Base:  TRY,
			Rates: map[CurrencyCode]float64{TRY: 1, USD: 8},
		}, nil).
		Once()

	service := NewService(mockProvider)

	res, err := service.FetchRates(context.Background(), &FetchParams{Base: TRY})
	assert.NoError(t, err)
	assert.Equal(t, &RatesResponse{Base: TRY, Rates: map[CurrencyCode]float64{TRY: 1, USD: 8}}, res)

	mockProvider.AssertExpectations(t)
}
