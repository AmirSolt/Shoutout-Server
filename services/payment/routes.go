package payment

import (
	"basedpocket/base"
	"fmt"
	"io"
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase/core"
	"github.com/stripe/stripe-go/v78/webhook"
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
		return ctx.String(http.StatusServiceUnavailable, fmt.Errorf("problem with request. Error: %w", err).Error())
	}
	event, err := webhook.ConstructEvent(payload, req.Header.Get("Stripe-Signature"), env.STRIPE_WEBHOOK_KEY)
	if err != nil {
		return ctx.String(http.StatusBadRequest, fmt.Errorf("error verifying webhook signature. Error: %w", err).Error())
	}
	// ==================================================================

	if err := onStripeEvents(app, ctx, event); err != nil {
		return ctx.String(http.StatusInternalServerError, err.Error())
	}

	res.Writer.WriteHeader(http.StatusOK)
	return nil
}
