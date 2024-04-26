package payment

import (
	"basedpocket/base"
	"fmt"
	"io"
	"net/http"

	"github.com/getsentry/sentry-go"
	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase/core"
	"github.com/stripe/stripe-go/v76/webhook"
)

func handleStripeWebhook(app core.App, ctx echo.Context, env *base.Env) error {
	// ==================================================================
	// The signature check is pulled directly from Stripe and it's not tested
	req := ctx.Request()
	res := ctx.Response()

	const MaxBodyBytes = int64(65536)
	req.Body = http.MaxBytesReader(res.Writer, req.Body, MaxBodyBytes)
	payload, err := io.ReadAll(req.Body)
	if err != nil {
		eventID := sentry.CaptureException(err)
		return ctx.String(http.StatusServiceUnavailable, fmt.Errorf("problem with request. eventID: %s", *eventID).Error())
	}
	endpointSecret := env.STRIPE_WEBHOOK_KEY
	event, err := webhook.ConstructEvent(payload, req.Header.Get("Stripe-Signature"), endpointSecret)
	if err != nil {
		eventID := sentry.CaptureException(err)
		return ctx.String(http.StatusBadRequest, fmt.Errorf("error verifying webhook signature. eventID: %s", *eventID).Error())
	}
	// ==================================================================

	if err := onStripeEvents(app, ctx, event); err != nil {
		return ctx.String(http.StatusInternalServerError, err.Error.Error())
	}

	res.Writer.WriteHeader(http.StatusOK)
	return nil
}
