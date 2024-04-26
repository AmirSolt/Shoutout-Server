package main

import (
	"basedpocket/base"
	"basedpocket/cmodels"
	"basedpocket/services/notif"
	"basedpocket/services/payment"
	"log"

	"github.com/pocketbase/pocketbase"
)

// go run main.go serve
//

func main() {

	env := base.LoadEnv()
	base.LoadLogging(env)
	app := pocketbase.New()

	cmodels.LoadModels(app, env)
	payment.LoadPayment(app, env)
	notif.LoadNotif(app, env)

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
