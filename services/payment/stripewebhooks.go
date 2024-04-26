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

	if event.Type == "customer.created" {
		return onCustomerCreatedEvent(app, ctx, event)
	}
	if event.Type == "customer.deleted" {
		return onCustomerDeletedEvent(app, ctx, event)
	}
	if event.Type == "customer.subscription.created" {
		return onSubscriptionCreatedEvent(app, ctx, event)
	}
	if event.Type == "customer.subscription.updated" {
		return onSubscriptionUpdatedEvent(app, ctx, event)
	}
	if event.Type == "customer.subscription.deleted" {
		return onSubscriptionDeletedEvent(app, ctx, event)
	}

	err := fmt.Errorf("unhandled stripe event type: %s\n", event.Type)
	eventID := sentry.CaptureException(err)
	return &base.CError{Message: "Internal Server Error", EventID: *eventID, Error: err}
}

// ===============================================================================

func onCustomerCreatedEvent(app core.App, ctx echo.Context, event stripe.Event) *base.CError {
	stripeCustomer, err := getStripeCustomerFromObj(event.Data.Object)
	if err != nil {
		return err
	}

	var user *cmodels.User
	if err := app.Dao().ModelQuery(user).
		AndWhere(dbx.HashExp{"email": stripeCustomer.Email}).
		Limit(1).
		One(&user); err != nil {
		return cmodels.HandleReadError(err, false)
	}

	newCustomer := &cmodels.Customer{
		User:                 user.Id,
		StripeCustomerID:     stripeCustomer.ID,
		StripeSubscriptionID: "",
		Tier:                 0,
	}
	if err := cmodels.Save(app, newCustomer); err != nil {
		return err
	}
	return nil
}

// ===============================================================================

func onCustomerDeletedEvent(app core.App, ctx echo.Context, event stripe.Event) *base.CError {
	stripeCustomer, err := getStripeCustomerFromObj(event.Data.Object)
	if err != nil {
		return err
	}

	var customer *cmodels.Customer
	if err := app.Dao().ModelQuery(customer).
		AndWhere(dbx.HashExp{"stripe_customer_id": stripeCustomer.ID}).
		Limit(1).
		One(&customer); err != nil {
		return cmodels.HandleReadError(err, false)
	}

	if err := cmodels.Delete(app, customer); err != nil {
		return err
	}
	return nil
}

// ===============================================================================

func onSubscriptionCreatedEvent(app core.App, ctx echo.Context, event stripe.Event) *base.CError {
	stripeSubscription, err := getStripeSubscriptionFromObj(event.Data.Object)
	if err != nil {
		return err
	}

	tier, err := getSubscriptionTier(stripeSubscription)
	if err != nil {
		return err
	}

	var customer *cmodels.Customer
	if err := app.Dao().ModelQuery(customer).
		AndWhere(dbx.HashExp{"stripe_customer_id": stripeSubscription.Customer.ID}).
		Limit(1).
		One(&customer); err != nil {
		return cmodels.HandleReadError(err, false)
	}

	customer.StripeSubscriptionID = stripeSubscription.ID
	customer.Tier = tier
	if err := cmodels.Save(app, customer); err != nil {
		return err
	}
	return nil
}

// ===============================================================================

func onSubscriptionUpdatedEvent(app core.App, ctx echo.Context, event stripe.Event) *base.CError {
	stripeSubscription, err := getStripeSubscriptionFromObj(event.Data.Object)
	if err != nil {
		return err
	}

	tier, err := getSubscriptionTier(stripeSubscription)
	if err != nil {
		return err
	}

	var customer *cmodels.Customer
	if err := app.Dao().ModelQuery(customer).
		AndWhere(dbx.HashExp{"stripe_customer_id": stripeSubscription.Customer.ID}).
		Limit(1).
		One(&customer); err != nil {
		return cmodels.HandleReadError(err, false)
	}

	customer.Tier = tier
	if err := cmodels.Save(app, customer); err != nil {
		return err
	}
	return nil
}

// ===============================================================================

func onSubscriptionDeletedEvent(app core.App, ctx echo.Context, event stripe.Event) *base.CError {
	stripeSubscription, err := getStripeSubscriptionFromObj(event.Data.Object)
	if err != nil {
		return err
	}

	var customer *cmodels.Customer
	if err := app.Dao().ModelQuery(customer).
		AndWhere(dbx.HashExp{"stripe_customer_id": stripeSubscription.Customer.ID}).
		Limit(1).
		One(&customer); err != nil {
		return cmodels.HandleReadError(err, false)
	}

	if stripeSubscription.ID == customer.StripeSubscriptionID {
		customer.StripeSubscriptionID = ""
		customer.Tier = 0
		if err := cmodels.Save(app, customer); err != nil {
			return err
		}
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
