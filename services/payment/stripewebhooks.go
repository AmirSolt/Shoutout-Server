package payment

import (
	"basedpocket/cmodels"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase/core"
	"github.com/stripe/stripe-go/v78"
)

func onStripeEvents(app core.App, ctx echo.Context, event stripe.Event) error {
	if event.Type == "checkout.session.completed" {
		return onCheckoutSuccess(app, ctx, event)
	}

	err := fmt.Errorf("unhandled stripe event type: %s", event.Type)
	return err
}

// ===============================================================================

func onCheckoutSuccess(app core.App, ctx echo.Context, event stripe.Event) error {
	checkoutSession, err := getStripeCheckoutSessionFromObj(event.Data.Object)
	if err != nil {
		return err
	}

	trans := &cmodels.Transaction{
		Amount:        float64(checkoutSession.AmountTotal) / float64(100),
		PaymentIntent: checkoutSession.PaymentIntent.ID,
		UserName:      checkoutSession.Metadata["user_name"],
		UserEmail:     checkoutSession.CustomerDetails.Email,
		CharacterName: checkoutSession.Metadata["character_name"],
	}

	if err := cmodels.Save(app, trans); err != nil {
		return err
	}
	return nil
}

// ===============================================================================
// ===============================================================================
// ===============================================================================

func getStripeCustomerFromObj(object map[string]interface{}) (*stripe.Customer, error) {
	jsonCustomer, err := json.Marshal(object)
	if err != nil {
		return nil, err
	}
	var stripeCustomer *stripe.Customer
	err = json.Unmarshal(jsonCustomer, &stripeCustomer)
	if stripeCustomer == nil || err != nil {
		return nil, err
	}
	return stripeCustomer, nil
}

func getStripeCheckoutSessionFromObj(object map[string]interface{}) (*stripe.CheckoutSession, error) {
	jsonCustomer, err := json.Marshal(object)
	if err != nil {

		return nil, err
	}
	var stripeStruct *stripe.CheckoutSession
	err = json.Unmarshal(jsonCustomer, &stripeStruct)
	if stripeStruct == nil || err != nil {

		return nil, err
	}
	return stripeStruct, nil
}

func getStripeSubscriptionFromObj(object map[string]interface{}) (*stripe.Subscription, error) {
	jsonCustomer, err := json.Marshal(object)
	if err != nil {

		return nil, err
	}
	var stripeStruct *stripe.Subscription
	err = json.Unmarshal(jsonCustomer, &stripeStruct)
	if stripeStruct == nil || err != nil {

		return nil, err
	}
	return stripeStruct, nil
}

func getSubscriptionTier(subsc *stripe.Subscription) (int, error) {
	if subsc == nil {
		return 0, nil
	}
	subscTierStr := subsc.Items.Data[0].Price.Metadata["tier"]
	subscTierInt, err := strconv.Atoi(subscTierStr)
	if err != nil {

		return 0, err
	}
	return subscTierInt, nil
}
