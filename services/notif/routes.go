package notif

import (
	"basedpocket/base"
	"basedpocket/cmodels"
	"fmt"
	"net/http"
	"net/mail"

	"github.com/getsentry/sentry-go"
	"github.com/labstack/echo/v5"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/mailer"
)

type NotifCreateManyParams struct {
	Params []NotifCreateParams `json:"params"`
}

type NotifCreateParams struct {
	UserID   string `db:"user_id" json:"user_id"`
	Subject  string `db:"subject" json:"subject"`
	BodyHTML string `db:"body_html" json:"body_html"`
}

func handleCreateNotifs(app core.App, ctx echo.Context, env *base.Env) error {

	var notifCreateManyParams NotifCreateManyParams
	err := ctx.Bind(&notifCreateManyParams)
	if err != nil {
		eventID := sentry.CaptureException(err)
		cerr := &base.CError{Message: "Internal Server Error", EventID: *eventID, Error: err}
		return ctx.JSON(http.StatusInternalServerError, cerr)
	}

	app.Logger().Info(fmt.Sprintf("handleCreateNotifs Params: %+v", notifCreateManyParams))

	for _, param := range notifCreateManyParams.Params {

		var user *cmodels.User
		if err := app.Dao().ModelQuery(user).
			AndWhere(dbx.HashExp{"id": param.UserID}).
			Limit(1).
			One(&user); err != nil {
			cmodels.HandleReadError(err, false)
			continue
		}

		success := sendEmail(app, user.Email, param)

		notif := &cmodels.Notif{
			User:             param.UserID,
			ToEmail:          user.Email,
			Subject:          param.Subject,
			BodyHTML:         param.BodyHTML,
			SendingAttempted: true,
			IsSuccessful:     success,
		}
		if cerr := cmodels.Save(app, notif); cerr != nil {
			continue
		}
	}

	return ctx.NoContent(http.StatusOK)
}

func sendEmail(app core.App, toEmail string, params NotifCreateParams) bool {
	message := &mailer.Message{
		From: mail.Address{
			Address: app.Settings().Meta.SenderAddress,
			Name:    app.Settings().Meta.SenderName,
		},
		To:      []mail.Address{{Address: toEmail}},
		Subject: params.Subject,
		HTML:    params.BodyHTML,
	}

	err := app.NewMailClient().Send(message)
	if err != nil {
		sentry.CaptureException(err)
		return false
	}

	return true
}
