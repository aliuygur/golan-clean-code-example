package main

import (
	"log"
	"net/http"
	"os"

	"github.com/alioygur/golang-clean-code-example/currencyrates"
	"github.com/alioygur/golang-clean-code-example/currencyrates/currencyratesapi"
	"github.com/alioygur/golang-clean-code-example/currencyrates/providers/tcmb"
	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	// init currencyrates service
	currencyratesProviderTCMB := tcmb.NewClient(nil)
	currencyratesService := currencyrates.NewService(currencyratesProviderTCMB)
	handler := &currencyratesapi.Handler{CurrencyRates: currencyratesService}
	handler.RegisterRoutes(r)

	// start the http server
	log.Printf("listening on port %s", getDefaultPort())
	if err := http.ListenAndServe(":"+getDefaultPort(), r); err != nil {
		log.Fatalf("unable to start http server: %s", err)
	}
}

func getDefaultPort() string {
	if p := os.Getenv("PORT"); p != "" {
		return p
	}
	return "8080"
}
