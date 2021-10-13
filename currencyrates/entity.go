package currencyrates

type RatesResponse struct {
	Base  CurrencyCode
	Rates map[CurrencyCode]float64
}

// CurrencyCode represents currency code type
type CurrencyCode string

// Currency Codes
const (
	TRY CurrencyCode = "TRY"
	USD CurrencyCode = "USD"
	AUD CurrencyCode = "AUD"
	DKK CurrencyCode = "DKK"
	EUR CurrencyCode = "EUR"
	GBP CurrencyCode = "GBP"
	CHF CurrencyCode = "CHF"
	SEK CurrencyCode = "SEK"
	CAD CurrencyCode = "CAD"
	KWD CurrencyCode = "KWD"
	NOK CurrencyCode = "NOK"
	SAR CurrencyCode = "SAR"
	JPY CurrencyCode = "JPY"
	BGN CurrencyCode = "BGN"
	RON CurrencyCode = "RON"
	RUB CurrencyCode = "RUB"
	IRR CurrencyCode = "IRR"
	CNY CurrencyCode = "CNY"
	PKR CurrencyCode = "PKR"
	QAR CurrencyCode = "QAR"
	KRW CurrencyCode = "KRW"
)

var AvailableCurrencies = []CurrencyCode{
	TRY,
	USD,
	AUD,
	DKK,
	EUR,
	GBP,
	CHF,
	SEK,
	CAD,
	KWD,
	NOK,
	SAR,
	JPY,
	BGN,
	RON,
	RUB,
	IRR,
	CNY,
	PKR,
	QAR,
	KRW,
}

func IsCodeAvailable(code CurrencyCode) bool {
	for _, c := range AvailableCurrencies {
		if code == c {
			return true
		}
	}
	return false
}
