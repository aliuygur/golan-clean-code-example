package currencyrates

import (
	"context"
)

type rateprovider interface {
	FetchRates(context.Context, *FetchParams) (*RatesResponse, error)
}

func NewService(rates rateprovider) *Service {
	return &Service{Rates: rates}
}

type Service struct {
	Rates rateprovider
}

type FetchParams struct {
	Base CurrencyCode
}

func (s *Service) FetchRates(ctx context.Context, params *FetchParams) (*RatesResponse, error) {
	if !IsCodeAvailable(params.Base) {
		return nil, ErrInvalidCurrencyCode
	}

	return s.Rates.FetchRates(ctx, params)
}

func (s *Service) Exchange(ctx context.Context, amount float64, from CurrencyCode, to CurrencyCode) (float64, error) {
	rate, err := s.GetRate(ctx, from, to)
	if err != nil {
		return 0, err
	}

	total := amount * rate

	return total, nil
}

func (s *Service) GetRate(ctx context.Context, from CurrencyCode, to CurrencyCode) (float64, error) {
	if !IsCodeAvailable(from) || !IsCodeAvailable(to) {
		return 0, ErrInvalidCurrencyCode
	}
	res, err := s.FetchRates(ctx, &FetchParams{Base: TRY})
	if err != nil {
		return 0, err
	}

	return res.Rates[to] / res.Rates[from], nil
}
