package cmodels

import (
	"database/sql"
	"errors"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/models"
)

func FindOneQuery[T any](app core.App, ctx echo.Context, item *T, queryStr string, params dbx.Params, skipNoRowsErr bool) error {
	// params["table"] = any(item).(models.Model).TableName()
	err := app.Dao().DB().NewQuery(queryStr).Bind(params).WithContext(ctx.Request().Context()).One(item)
	return HandleReadError(err, skipNoRowsErr)
}

func FindAllQuery[T any](app core.App, ctx echo.Context, items *[]T, queryStr string, params dbx.Params, skipNoRowsErr bool) error {
	// params["table"] = any(item).(models.Model).TableName()
	err := app.Dao().DB().NewQuery(queryStr).Bind(params).WithContext(ctx.Request().Context()).All(items)
	return HandleReadError(err, skipNoRowsErr)
}

func HandleReadError(err error, skipNoRowsErr bool) error {
	if err != nil && !errors.Is(err, sql.ErrNoRows) {

		return err
	}
	if errors.Is(err, sql.ErrNoRows) && !skipNoRowsErr {

		return err
	}
	if errors.Is(err, sql.ErrNoRows) && skipNoRowsErr {
		return nil
	}
	return nil
}

func Save[T any](app core.App, item T) error {
	if err := app.Dao().Save(any(item).(models.Model)); err != nil {

		return err
	}
	return nil
}
func Delete[T any](app core.App, item T) error {
	if err := app.Dao().Delete(any(item).(models.Model)); err != nil {

		return err
	}
	return nil
}
