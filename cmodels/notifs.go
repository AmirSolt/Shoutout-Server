package cmodels

import (
	"fmt"
	"log"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/models/schema"
	"github.com/pocketbase/pocketbase/tools/types"
)

const notifs string = "notifs"

var _ models.Model = (*Notif)(nil)

type Notif struct {
	models.BaseModel
	User             string `db:"user" json:"user"`
	ToEmail          string `db:"to_email" json:"to_email"`
	Subject          string `db:"subject" json:"subject"`
	BodyHTML         string `db:"body_html" json:"body_html"`
	SendingAttempted bool   `db:"sending_attempted" json:"sending_attempted"`
	IsSuccessful     bool   `db:"is_successful" json:"is_successful"`
}

func (m *Notif) TableName() string {
	return notifs // the name of your collection
}

// =======================================

func createNotifsCollection(app core.App) {

	collectionName := notifs

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
		ListRule:   nil,
		ViewRule:   nil,
		CreateRule: nil,
		UpdateRule: nil,
		DeleteRule: nil,
		Schema: schema.NewSchema(
			&schema.SchemaField{
				Name:     "user",
				Type:     schema.FieldTypeRelation,
				Required: true,
				Options: &schema.RelationOptions{
					CollectionId:  users.Id,
					CascadeDelete: true,
				},
			},
			&schema.SchemaField{
				Name:     "to_email",
				Type:     schema.FieldTypeEmail,
				Required: true,
				Options:  &schema.EmailOptions{},
			},
			&schema.SchemaField{
				Name:     "subject",
				Type:     schema.FieldTypeText,
				Required: true,
				Options:  &schema.TextOptions{},
			},
			&schema.SchemaField{
				Name:     "body_html",
				Type:     schema.FieldTypeText,
				Required: true,
				Options:  &schema.TextOptions{},
			},
			&schema.SchemaField{
				Name:     "sending_attempted",
				Type:     schema.FieldTypeBool,
				Required: true,
				Options:  &schema.BoolOptions{},
			},
			&schema.SchemaField{
				Name:     "is_successful",
				Type:     schema.FieldTypeBool,
				Required: true,
				Options:  &schema.BoolOptions{},
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
