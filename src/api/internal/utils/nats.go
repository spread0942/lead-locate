package utils

import (
	"time"
	"encoding/json"

	"github.com/nats-io/nats.go"

	"api/internal/utils/natsUtils"
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
	return natsCtl, nil
}

// handle the request and response
func (nats *NatsCtl) Request(subj string, method string, body any) (any, error) {
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
	reply, err := nats.client.Request(subj, data, nats.Timeout)
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