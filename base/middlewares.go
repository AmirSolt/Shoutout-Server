package base

import (
	"net/http"

	"github.com/labstack/echo/v5"
)

const HeaderSecretKeyName string = "Api-Key"

func GeoServerSecretMiddleware(env *Env) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {

			expectedString := map[string]bool{
				env.PRIVATE_GEO_SERVER_API_KEY: true,
			}

			for key, values := range ctx.Request().Header {
				if key == HeaderSecretKeyName {
					for _, value := range values {
						if expectedString[value] {
							next(ctx)
						}
					}
				}
			}
			return ctx.NoContent(http.StatusUnauthorized)

		}
	}
}
