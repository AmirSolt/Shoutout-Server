package cmodels

import (
	"fmt"
	"log"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/models/schema"
	"github.com/pocketbase/pocketbase/tools/types"
)

const customers string = "customers"

var _ models.Model = (*Customer)(nil)

type Customer struct {
	models.BaseModel
	User                 string `db:"user" json:"user"`
	StripeCustomerID     string `db:"stripe_customer_id" json:"stripe_customer_id"`
	StripeSubscriptionID string `db:"stripe_subscription_id" json:"stripe_subscription_id"`
	Tier                 int    `db:"tier" json:"tier"`
}

func (m *Customer) TableName() string {
	return customers // the name of your collection
}

// =======================================

// func (c *Customer) Find(app core.App, ctx echo.Context, skipNoRowsErr bool) *base.CError {
// 	return FindOne(app, ctx, c, "")
// }

// func (c *Customer) Save(app core.App, ctx echo.Context) *base.CError {

// }

// =======================================

func createCustomersCollection(app core.App) {

	collectionName := customers

	existingCollection, _ := app.Dao().FindCollectionByNameOrId(collectionName)
	if existingCollection != nil {
		return
	}

	users, err := app.Dao().FindCollectionByNameOrId(users)
	if err != nil {
		log.Fatalf("users table not found: %+v", err)
	}

	collection := &models.Collection{
		Name:       collectionName,
		Type:       models.CollectionTypeBase,
		ListRule:   types.Pointer("user.id = @request.auth.id"),
		ViewRule:   types.Pointer("user.id = @request.auth.id"),
		CreateRule: nil,
		UpdateRule: nil,
		DeleteRule: nil,
		Schema: schema.NewSchema(
			&schema.SchemaField{
				Name:     "user",
				Type:     schema.FieldTypeRelation,
				Required: true,
				Options: &schema.RelationOptions{
					MaxSelect:     types.Pointer(1),
					CollectionId:  users.Id,
					CascadeDelete: true,
				},
			},
			&schema.SchemaField{
				Name:     "stripe_customer_id",
				Type:     schema.FieldTypeText,
				Required: true,
				Options:  &schema.TextOptions{},
			},
			&schema.SchemaField{
				Name:     "stripe_subscription_id",
				Type:     schema.FieldTypeText,
				Required: false,
				Options:  &schema.TextOptions{},
			},
			&schema.SchemaField{
				Name:     "tier",
				Type:     schema.FieldTypeNumber,
				Required: true,
				Options:  &schema.NumberOptions{},
			},
		),
		Indexes: types.JsonArray[string]{
			fmt.Sprintf("CREATE UNIQUE INDEX idx_user ON %s (user)", collectionName),
		},
	}

	if err := app.Dao().SaveCollection(collection); err != nil {
		log.Fatalf("%s collection failed: %+v", collectionName, err)
	}
}
