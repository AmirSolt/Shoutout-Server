package cmodels

import (
	"fmt"
	"log"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/models/schema"
)

const orders string = "orders"

var _ models.Model = (*Order)(nil)

type Order struct {
	models.BaseModel
	PaymentIntent string `db:"payment_intent" json:"payment_intent"`
	Status        string `db:"status" json:"status"`
	Images        string `db:"images" json:"images"`
	Message       string `db:"message" json:"message"`
}

func (m *Order) TableName() string {
	return orders // the name of your collection
}

// =======================================

type OrderStatus string

const (
	PaymentPending OrderStatus = "payment_pending"
	PaymentFailed  OrderStatus = "payment_failed"
	OrderWaiting   OrderStatus = "order_waiting"
	OrderComplete  OrderStatus = "order_complete"
	OrderRejected  OrderStatus = "order_rejected"
	OrderCancelled OrderStatus = "order_cancelled"
)

func getOrderStatusRegex() string {
	return fmt.Sprintf(`^(%s|%s|%s|%s|%s|%s)$`, PaymentPending, PaymentFailed, OrderWaiting, OrderComplete, OrderRejected, OrderCancelled)
}

// =======================================

func createOrdersCollection(app core.App) {

	collectionName := orders

	existingCollection, _ := app.Dao().FindCollectionByNameOrId(collectionName)
	if existingCollection != nil {
		return
	}

	collection := &models.Collection{
		Name:       collectionName,
		Type:       models.CollectionTypeBase,
		ListRule:   nil,
		ViewRule:   nil,
		CreateRule: nil,
		UpdateRule: nil,
		DeleteRule: nil,
		Schema: schema.NewSchema(
			&schema.SchemaField{
				Name:     "payment_intent",
				Type:     schema.FieldTypeText,
				Required: false,
				Options:  &schema.TextOptions{},
			},
			&schema.SchemaField{
				Name:     "status",
				Type:     schema.FieldTypeText,
				Required: true,
				Options:  &schema.TextOptions{Pattern: getOrderStatusRegex()},
			},
			&schema.SchemaField{
				Name:     "images",
				Type:     schema.FieldTypeFile,
				Required: true,
				Options:  &schema.FileOptions{Protected: false, MaxSelect: 10, MaxSize: 1_000_000, MimeTypes: []string{".jpg", ".png", ".jpeg"}},
			},
			&schema.SchemaField{
				Name:     "message",
				Type:     schema.FieldTypeText,
				Required: true,
				Options:  &schema.TextOptions{},
			},
		),
	}

	if err := app.Dao().SaveCollection(collection); err != nil {
		log.Fatalf("%s collection failed: %+v", collectionName, err)
	}
}
