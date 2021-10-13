package app

import (
	"context"
	"log"
	"net/http"

	"cloud.google.com/go/compute/metadata"
	"cloud.google.com/go/errorreporting"
)

var gcperr *errorreporting.Client

func initErrorReportingClient() {
	if pid, err := metadata.ProjectID(); err == nil && pid != "" {
		ec, err := errorreporting.NewClient(context.Background(), pid, errorreporting.Config{
			ServiceName:    mustGetEnv("K_SERVICE"),
			ServiceVersion: mustGetEnv("K_REVISION"),
			OnError: func(err error) {
				log.Printf("could not log error: %v", err)
			},
		})
		if err != nil {
			log.Fatal(err)
		}
		gcperr = ec
	}
}

func ReportError(r *http.Request, err error) {
	if gcperr != nil {
		gcperr.Report(errorreporting.Entry{
			Error: err,
			Req:   r,
		})
	}
}
