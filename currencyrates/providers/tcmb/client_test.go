package tcmb

import (
	"context"
	"encoding/xml"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/alioygur/golang-clean-code-example/currencyrates"
)

type roundTripperFunc func(*http.Request) (*http.Response, error)

func (f roundTripperFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func Test_decodeAPIResponse(t *testing.T) {
	data, err := os.ReadFile("./testdata/tcmb_rates_response.xml")
	require.NoError(t, err)

	var res ratesResponse
	require.NoError(t, xml.Unmarshal(data, &res))
}

func TestClient_FetchRates(t *testing.T) {
	testTransport := roundTripperFunc(func(r *http.Request) (*http.Response, error) {
		f, err := os.Open("testdata/tcmb_rates_response.xml")
		require.NoError(t, err)
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       f,
			Header:     make(http.Header),
		}, nil
	})

	client := NewClient(testTransport)

	opts := &currencyrates.FetchParams{Base: currencyrates.TRY}
	ratesRes, err := client.FetchRates(context.Background(), opts)
	require.NoError(t, err)

	assert.Equal(t, opts.Base, ratesRes.Base)
	assert.Equal(t, float64(1), ratesRes.Rates[opts.Base])
	assert.Equal(t, 21, len(ratesRes.Rates))
	assert.Equal(t, 1.0, ratesRes.Rates[currencyrates.TRY])
	assert.Equal(t, 8.8399, ratesRes.Rates[currencyrates.USD])
	assert.Equal(t, 0.080213, ratesRes.Rates[currencyrates.JPY])
}
