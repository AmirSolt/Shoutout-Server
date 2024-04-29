package base

import (
	"github.com/getsentry/sentry-go"
)

type CError struct {
	Error   error          `json:"-"`
	Message string         `json:"message"`
	EventID sentry.EventID `json:"eventID"`
}
