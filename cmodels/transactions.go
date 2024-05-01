package cmodels

import (
	"log"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/models/schema"
	"go.openly.dev/pointy"
)

const transactions string = "transactions"

var _ models.Model = (*Transaction)(nil)

type Transaction struct {
	models.BaseModel
	Amount        float64 `db:"amount" json:"amount"`
	PaymentIntent string  `db:"payment_intent" json:"payment_intent"`
	UserName      string  `db:"user_name" json:"user_name"`
	CharacterName string  `db:"character_name" json:"character_name"`
}

func (m *Transaction) TableName() string {
	return transactions // the name of your collection
}

// =======================================

func createTransactionsCollection(app core.App) {

	collectionName := transactions

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
				Name:     "amount",
				Type:     schema.FieldTypeNumber,
				Required: true,
				Options:  &schema.NumberOptions{Min: pointy.Float64(4.99), Max: pointy.Float64(10_001.0)},
			},
			&schema.SchemaField{
				Name:     "payment_intent",
				Type:     schema.FieldTypeText,
				Required: true,
				Options:  &schema.TextOptions{},
			},
			&schema.SchemaField{
				Name:     "user_name",
				Type:     schema.FieldTypeText,
				Required: true,
				Options:  &schema.TextOptions{Min: pointy.Int(3), Max: pointy.Int(24)},
			},
			&schema.SchemaField{
				Name:     "character_name",
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
