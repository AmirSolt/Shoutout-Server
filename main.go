package main

import (
	"basedpocket/base"
	"basedpocket/cmodels"
	"basedpocket/services/payment"
	"log"

	"github.com/pocketbase/pocketbase"
)

// go run main.go serve
//

func main() {

	env := base.LoadEnv()
	app := pocketbase.New()

	cmodels.LoadModels(app, env)
	payment.LoadPayment(app, env)

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
