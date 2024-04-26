package notif

import (
	"basedpocket/base"
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
)

func LoadNotif(app *pocketbase.PocketBase, env *base.Env) {

	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		// ===================
		// routes
		e.Router.AddRoute(echo.Route{
			Method: http.MethodPost,
			Path:   "/api/notifs/create/many",
			Handler: func(c echo.Context) error {
				return handleCreateNotifs(e.App, c, env)
			},
			Middlewares: []echo.MiddlewareFunc{
				apis.ActivityLogger(e.App),
				base.GeoServerSecretMiddleware(env),
			},
		})

		return nil
	})
}
