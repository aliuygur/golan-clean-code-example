package app

import (
	"os"
	"strings"
)

// env
const (
	DEV = iota
	TEST
	STAGE
	PROD
)

var ENV = DEV
var BaseURL = "http://localhost:8080"

func Init(env, baseurl string) {
	env = strings.ToLower(env)
	switch env {
	case "test":
		ENV = TEST
	case "stage":
		ENV = STAGE
	case "prod", "production":
		ENV = PROD
	default:
		ENV = DEV
	}

	BaseURL = strings.TrimRight(baseurl, "/")

	initErrorReportingClient()
}

func IsDEV() bool {
	return ENV == DEV
}

func mustGetEnv(key string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	panic("env required: " + key)
}
