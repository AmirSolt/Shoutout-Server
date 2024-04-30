package payment

import (
	"basedpocket/base"
	"basedpocket/cmodels"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/getsentry/sentry-go"
	"github.com/labstack/echo/v5"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
	"github.com/stripe/stripe-go/v76"
)

func onStripeEvents(app core.App, ctx echo.Context, event stripe.Event) *base.CError {
	if event.Type == "checkout.session.async_payment_succeeded" {
		return onCheckoutSuccess(app, ctx, event)
	}
	if event.Type == "checkout.session.async_payment_failed" {
		return onCheckoutFail(app, ctx, event)
	}

	err := fmt.Errorf("unhandled stripe event type: %s\n", event.Type)
	eventID := sentry.CaptureException(err)
	return &base.CError{Message: "Internal Server Error", EventID: *eventID, Error: err}
}

// ===============================================================================

func onCheckoutSuccess(app core.App, ctx echo.Context, event stripe.Event) *base.CError {
	checkoutSession, err := getStripeCheckoutSessionFromObj(event.Data.Object)
	if err != nil {
		return err
	}

	var order *cmodels.Order
	if err := app.Dao().ModelQuery(order).
		AndWhere(dbx.HashExp{"checkout_id": checkoutSession.ID}).
		Limit(1).
		One(&order); err != nil {
		return cmodels.HandleReadError(err, false)
	}

	order.PaymentIntent = checkoutSession.PaymentIntent.ID
	order.Status = string(cmodels.OrderWaiting)
	if err := cmodels.Save(app, order); err != nil {
		return err
	}
	return nil
}

// ===============================================================================

func onCheckoutFail(app core.App, ctx echo.Context, event stripe.Event) *base.CError {
	checkoutSession, err := getStripeCheckoutSessionFromObj(event.Data.Object)
	if err != nil {
		return err
	}

	orderID := checkoutSession.Metadata["order_id"]

	var order *cmodels.Order
	if err := app.Dao().ModelQuery(order).
		AndWhere(dbx.HashExp{"id": orderID}).
		Limit(1).
		One(&order); err != nil {
		return cmodels.HandleReadError(err, false)
	}

	order.PaymentIntent = checkoutSession.PaymentIntent.ID
	order.Status = string(cmodels.PaymentFailed)
	if err := cmodels.Save(app, order); err != nil {
		return err
	}
	return nil
}

// ===============================================================================
// ===============================================================================
// ===============================================================================

func getStripeCustomerFromObj(object map[string]interface{}) (*stripe.Customer, *base.CError) {
	jsonCustomer, err := json.Marshal(object)
	if err != nil {
		eventID := sentry.CaptureException(err)
		return nil, &base.CError{Message: "Internal Server Error", EventID: *eventID, Error: err}
	}
	var stripeCustomer *stripe.Customer
	err = json.Unmarshal(jsonCustomer, &stripeCustomer)
	if stripeCustomer == nil || err != nil {
		eventID := sentry.CaptureException(err)
		return nil, &base.CError{Message: "Internal Server Error", EventID: *eventID, Error: err}
	}
	return stripeCustomer, nil
}

func getStripeCheckoutSessionFromObj(object map[string]interface{}) (*stripe.CheckoutSession, *base.CError) {
	jsonCustomer, err := json.Marshal(object)
	if err != nil {
		eventID := sentry.CaptureException(err)
		return nil, &base.CError{Message: "Internal Server Error", EventID: *eventID, Error: err}
	}
	var stripeStruct *stripe.CheckoutSession
	err = json.Unmarshal(jsonCustomer, &stripeStruct)
	if stripeStruct == nil || err != nil {
		eventID := sentry.CaptureException(err)
		return nil, &base.CError{Message: "Internal Server Error", EventID: *eventID, Error: err}
	}
	return stripeStruct, nil
}

func getStripeSubscriptionFromObj(object map[string]interface{}) (*stripe.Subscription, *base.CError) {
	jsonCustomer, err := json.Marshal(object)
	if err != nil {
		eventID := sentry.CaptureException(err)
		return nil, &base.CError{Message: "Internal Server Error", EventID: *eventID, Error: err}
	}
	var stripeStruct *stripe.Subscription
	err = json.Unmarshal(jsonCustomer, &stripeStruct)
	if stripeStruct == nil || err != nil {
		eventID := sentry.CaptureException(err)
		return nil, &base.CError{Message: "Internal Server Error", EventID: *eventID, Error: err}
	}
	return stripeStruct, nil
}

func getSubscriptionTier(subsc *stripe.Subscription) (int, *base.CError) {
	if subsc == nil {
		return 0, nil
	}
	subscTierStr := subsc.Items.Data[0].Price.Metadata["tier"]
	subscTierInt, err := strconv.Atoi(subscTierStr)
	if err != nil {
		eventID := sentry.CaptureException(err)
		return 0, &base.CError{Message: "Internal Server Error", EventID: *eventID, Error: err}
	}
	return subscTierInt, nil
}
