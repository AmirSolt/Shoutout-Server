package base

import (
	"log"

	"github.com/getsentry/sentry-go"
)

type CError struct {
	Error   error          `json:"-"`
	Message string         `json:"message"`
	EventID sentry.EventID `json:"eventID"`
}

func LoadLogging(env *Env) {
	err := sentry.Init(sentry.ClientOptions{
		Dsn:              env.GLITCHTIP_DSN,
		EnableTracing:    true,
		TracesSampleRate: 1.0,
		Debug:            !env.IS_PROD,
	})
	if err != nil {
		log.Fatalf("sentry.Init: %s", err)
	}

	// base.Engine.Use(SentryGinNew(SentryGinOptions{}))

}
