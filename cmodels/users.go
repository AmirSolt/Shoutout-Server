package cmodels

import (
	"basedpocket/base"
	"fmt"

	"github.com/getsentry/sentry-go"
	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/models"
)

const users string = "users"

var _ models.Model = (*User)(nil)

type User struct {
	models.BaseModel
	Email string `db:"email" json:"email"`
}
type FindUserParams struct {
	Id    string `db:"id"`
	Email string `db:"email"`
}

func (m *User) TableName() string {
	return users // the name of your collection
}

// ===================================

func (user *User) GetUserByContext(ctx echo.Context) *base.CError {
	user = ctx.Get(apis.ContextAuthRecordKey).(*User)
	if user == nil {
		err := fmt.Errorf("user not found")
		eventID := sentry.CaptureException(err)
		return &base.CError{Message: "Internal Server Error", EventID: *eventID, Error: err}
	}
	return nil
}
