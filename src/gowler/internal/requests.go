package internal

import (
	"gowler/internal/requests"

	"fmt"

	"github.com/nats-io/nats.go"

	"gowler/internal/utils"
)

func RequestsHandler(nc *utils.NatsCtl, serviceName string) {
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