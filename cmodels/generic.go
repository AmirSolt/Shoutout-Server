package cmodels

import (
	"basedpocket/base"
	"database/sql"
	"errors"

	"github.com/getsentry/sentry-go"
	"github.com/labstack/echo/v5"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/models"
)

func FindOneQuery[T any](app core.App, ctx echo.Context, item *T, queryStr string, params dbx.Params, skipNoRowsErr bool) *base.CError {
	// params["table"] = any(item).(models.Model).TableName()
	err := app.Dao().DB().NewQuery(queryStr).Bind(params).WithContext(ctx.Request().Context()).One(item)
	return HandleReadError(err, skipNoRowsErr)
}

func FindAllQuery[T any](app core.App, ctx echo.Context, items *[]T, queryStr string, params dbx.Params, skipNoRowsErr bool) *base.CError {
	// params["table"] = any(item).(models.Model).TableName()
	err := app.Dao().DB().NewQuery(queryStr).Bind(params).WithContext(ctx.Request().Context()).All(items)
	return HandleReadError(err, skipNoRowsErr)
}

func HandleReadError(err error, skipNoRowsErr bool) *base.CError {
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		eventID := sentry.CaptureException(err)
		return &base.CError{Message: "Internal Server Error", EventID: *eventID, Error: err}
	}
	if errors.Is(err, sql.ErrNoRows) && !skipNoRowsErr {
		eventID := sentry.CaptureException(err)
		return &base.CError{Message: "Internal Server Error", EventID: *eventID, Error: err}
	}
	if errors.Is(err, sql.ErrNoRows) && skipNoRowsErr {
		return nil
	}
	return nil
}

func Save[T any](app core.App, item T) *base.CError {
	if err := app.Dao().Save(any(item).(models.Model)); err != nil {
		eventID := sentry.CaptureException(err)
		return &base.CError{Message: "Internal Server Error", EventID: *eventID, Error: err}
	}
	return nil
}
func Delete[T any](app core.App, item T) *base.CError {
	if err := app.Dao().Delete(any(item).(models.Model)); err != nil {
		eventID := sentry.CaptureException(err)
		return &base.CError{Message: "Internal Server Error", EventID: *eventID, Error: err}
	}
	return nil
}
