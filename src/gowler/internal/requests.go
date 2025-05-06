package internal

import (
	"gowler/internal/requests"

	"fmt"

	"github.com/nats-io/nats.go"
)

func RequestsHandler(nc *nats.Conn, serviceName string) {
	request := serviceName + ".request"
	mapsRequest := request + ".maps"
	nc.Subscribe(mapsRequest, func(msg *nats.Msg) {
		fmt.Println("[INFO] Received message:", string(msg.Data))
		requests.GetMaps(msg)
	})

	gowlerRequest := request + ".gowler"
	nc.Subscribe(gowlerRequest, func(msg *nats.Msg) {
		fmt.Println("[INFO] Received message:", string(msg.Data))
		requests.GowlerIt(msg)
	})
}