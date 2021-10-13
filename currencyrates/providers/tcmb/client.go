package tcmb

import (
	"context"
	"encoding/xml"
	"io/ioutil"
	"net/http"

	"github.com/alioygur/golang-clean-code-example/app"
	"github.com/alioygur/golang-clean-code-example/currencyrates"
)

func NewClient(rt http.RoundTripper) *Client {
	return &Client{client: &http.Client{Transport: rt}}
}

type Client struct {
	client *http.Client
}

type ratesResponse struct {
	XMLName  xml.Name `xml:"Tarih_Date"`
	Text     string   `xml:",chardata"`
	Tarih    string   `xml:"Tarih,attr"`
	Date     string   `xml:"Date,attr"`
	BultenNo string   `xml:"Bulten_No,attr"`
	Currency []struct {
		Text            string  `xml:",chardata"`
		CrossOrder      string  `xml:"CrossOrder,attr"`
		Kod             string  `xml:"Kod,attr"`
		CurrencyCode    string  `xml:"CurrencyCode,attr"`
		Unit            float64 `xml:"Unit"`
		Isim            string  `xml:"Isim"`
		CurrencyName    string  `xml:"CurrencyName"`
		ForexBuying     float64 `xml:"ForexBuying"`
		ForexSelling    float64 `xml:"ForexSelling"`
		BanknoteBuying  float64 `xml:"BanknoteBuying"`
		BanknoteSelling string  `xml:"BanknoteSelling"`
		CrossRateUSD    string  `xml:"CrossRateUSD"`
		CrossRateOther  string  `xml:"CrossRateOther"`
	} `xml:"Currency"`
}

func decodeAPIResponse(data []byte, v interface{}) error {
	return xml.Unmarshal(data, v)
}

func (tcmb *Client) FetchRates(ctx context.Context, opts *currencyrates.FetchParams) (*currencyrates.RatesResponse, error) {
	var todayURL = "https://www.tcmb.gov.tr/kurlar/today.xml"

	res, err := tcmb.client.Get(todayURL)
	if err != nil {
		return nil, app.Errorf("tcmb: unable to make http request: %w", err)
	}
	// to avoid resource leaks, call the res.body.close() method
	// https://stackoverflow.com/questions/33238518/what-could-happen-if-i-dont-close-response-body
	defer res.Body.Close()

	// read the request's body, for the future use.
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, app.Errorf("tcmb: unable to read response body: %w", err)
	}

	if res.StatusCode != http.StatusOK {
		return nil, app.Errorf("tcmb: unexpected response; status code: %d, body: %s", res.StatusCode, body)
	}

	var tcmbRes ratesResponse
	if err := decodeAPIResponse(body, &tcmbRes); err != nil {
		return nil, app.Errorf("tcmb: unable to decode response: %s", body)
	}

	ratesResponse := currencyrates.RatesResponse{
		Base:  currencyrates.TRY,
		Rates: make(map[currencyrates.CurrencyCode]float64),
	}

	for _, cr := range tcmbRes.Currency {
		// skip unsupported currency codes
		if !currencyrates.IsCodeAvailable(currencyrates.CurrencyCode(cr.CurrencyCode)) {
			continue
		}
		ratesResponse.Rates[currencyrates.CurrencyCode(cr.CurrencyCode)] = cr.ForexSelling / cr.Unit
	}
	ratesResponse.Rates[currencyrates.TRY] = 1

	return &ratesResponse, nil
}
