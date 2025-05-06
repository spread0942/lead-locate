package main

import (
	"log"
	"os"

	"api/internal"
)

const (
	serviceName = "api"
)

func main() {
	app := internal.NewApp(serviceName, "1.0.0")
	defer app.Stop()

	natsUrl := os.Getenv("NATS_URL")
	err := app.SetNatsContext(natsUrl)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%s >> Connected to NATS at %s", serviceName, natsUrl)

	apiUrl := os.Getenv("API_URL")
	app.StartApiServer(apiUrl)
	log.Printf("%s >> API server started at %s", serviceName, apiUrl)

	<-app.SrvCtx.Done()
	app.Stop()
}