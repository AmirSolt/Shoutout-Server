package pages

import (
	"basedpocket/base"
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
)

func LoadPages(app *pocketbase.PocketBase, env *base.Env) {

	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		// ===================
		// routes
		e.Router.AddRoute(echo.Route{
			Method: http.MethodGet,
			Path:   "/api/landing-page",
			Handler: func(c echo.Context) error {
				return handleLandingPageData(e.App, c, env)
			},
			Middlewares: []echo.MiddlewareFunc{
				apis.ActivityLogger(e.App),
			},
		})

		return nil
	})

}
