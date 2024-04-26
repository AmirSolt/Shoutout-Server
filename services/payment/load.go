package payment

import (
	"basedpocket/base"
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/stripe/stripe-go/v76"
)

func LoadPayment(app *pocketbase.PocketBase, env *base.Env) {

	stripe.Key = env.STRIPE_PRIVATE_KEY

	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		// ===================
		// routes
		e.Router.AddRoute(echo.Route{
			Method: http.MethodPost,
			Path:   "/api/stripe/webhook",
			Handler: func(c echo.Context) error {
				return handleStripeWebhook(e.App, c, env)
			},
			Middlewares: []echo.MiddlewareFunc{
				apis.ActivityLogger(e.App),
			},
		})

		return nil
	})

}
