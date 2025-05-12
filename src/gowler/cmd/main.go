package main

import (
	"log"
	"os"

	"gowler/internal"
)

const (
	serviceName = "gowler"
)

func main() {
	app := internal.NewApp(serviceName, "1.0.0")
	defer app.Stop()

	// nats configuration
	natsUrl := os.Getenv("NATS_URL")
	err := app.SetNatsContext(natsUrl)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%s >> Connected to NATS at %s", serviceName, natsUrl)

	go app.StartNatsSubscibes()

	<-app.SrvCtx.Done()
	app.Stop()
}
