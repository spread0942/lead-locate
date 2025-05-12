package utils

import (
	"time"
	"encoding/json"
	"log"

	"github.com/nats-io/nats.go"

	"gowler/internal/utils/natsUtils"
)

type NatsCtl struct {
	Url string
	Timeout time.Duration
	client *nats.Conn
}

// create an instace of natsCtl
// and connect on nats url
func NewNatsCtl(url string) (*NatsCtl, error) {
	nats, err := nats.Connect(url)
	if (err != nil) {
		return nil, err
	}
	natsCtl:= &NatsCtl{
		Url: url,
		Timeout: 10*time.Minute,
		client: nats,
	}

	// logs configuration
	initLog()

	return natsCtl, nil
}

// set up some log configuration
func initLog() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	log.Println("Logger initialized")
}

// handle the request and response
func (natsCtl *NatsCtl) Request(subj string, method string, body any) (any, error) {
	// prepare the request
	req := natsUtils.NatsRequest{
		Method: method,
		Timestamp: time.Now(),
		Body: body,
	}
	data, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	reply, err := natsCtl.client.Request(subj, data, natsCtl.Timeout)
	if err != nil {
		return nil, err
	}
	// elaborate the response
	resp := &natsUtils.NatsResponse{}
	err = json.Unmarshal(reply.Data, resp)
	if err != nil {
		return nil, err
	}
	if resp.Status != 0 {
		return nil, resp.Error
	}
	return resp.Body, nil
}

func (natsCtl *NatsCtl) SubscribeHandler(serviceName string) {
	subj := serviceName + ".request.maps"
	natsCtl.client.Subscribe(subj, func(msg *nats.Msg) {
		logSubscribe(subj, msg)

		var request natsUtils.NatsRequest
		if err := json.Unmarshal(msg.Data, &request); err != nil {
			log.Println("Error unmarshalling message:", err)
			return
		}

		rawBody, ok := request.Body.(json.RawMessage)
		if !ok {
			log.Println("Error: request.Body is not a valid JSON type")
            return
		}

		var body natsUtils.GowlerBody
		if err := json.Unmarshal(rawBody, &body); err != nil {
			log.Println("Error unmarshalling request.Body into GowlerBody:", err)
            return
		}

		requests.GetMaps(msg)
	})

	subj = serviceName + ".request.gowler"
	natsCtl.client.Subscribe(subj, func(msg *nats.Msg) {
		logSubscribe(subj, msg)
		requests.GowlerIt(msg)
	})
}

func logSubscribe(subj string, msg *nats.Msg) {
	log.Println("Subscribe:", subj)
	log.Println("Received message:", string(msg.Data))
}