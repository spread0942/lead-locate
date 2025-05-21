package natsCtl

import (
	"encoding/json"
	"log"
	"time"

	"github.com/nats-io/nats.go"

	"gowler/internal/google/maps"
	"gowler/internal/gowler"
	"gowler/internal/nominatim"
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
	req := NatsRequest{
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
	resp := &NatsResponse{}
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
	subj := serviceName + ".request.gowler"
	natsCtl.client.Subscribe(subj, func(msg *nats.Msg) {
		logSubscribe(subj, msg)
		reply := NatsResponse{
			Request: subj,
			Method: "gowler",
			Timestamp: time.Now(),
			Status: 0,
			Body: nil,
			Error: nil,
		}

		// Decode the request
		var request NatsRequest
		if err := json.Unmarshal(msg.Data, &request); err != nil {
			handleError(err, msg, &reply)
			return
		}

		var body GowlerBody
		rawBody, err := json.Marshal(request.Body)
		if err != nil {
			handleError(err, msg, &reply)
			return
		}
		if err := json.Unmarshal(rawBody, &body); err != nil {
			handleError(err, msg, &reply)
            return
		}

		// Crawl the website
		gogo := gowler.NewGowler(body.Website)
		gogo.Crawl()
		bodyReplay := GowlerReplyBody{}
		bodyReplay.Site = gogo.Site
		bodyReplay.Domain = gogo.Domain
		bodyReplay.SiteUrls = gogo.SiteUrls
		bodyReplay.OtherUrls = gogo.OtherUrls
		bodyReplay.Telephones = gogo.Telephones
		bodyReplay.Emails = gogo.Emails

		reply.Body = bodyReplay

		// Marshal the response
		responseData, err := json.Marshal(reply)
		if err != nil {
			handleError(err, msg, &reply)
			return
		}
		err = msg.Respond(responseData)
		if err != nil {
			handleError(err, msg, &reply)
			return
		}
	})

	subj = serviceName + ".request.maps"
	natsCtl.client.Subscribe(subj, func(msg *nats.Msg) {
		logSubscribe(subj, msg)
		reply := NatsResponse{
			Request: subj,
			Method: "gowler",
			Timestamp: time.Now(),
			Status: 0,
			Body: nil,
			Error: nil,
		}
		// Decode the request
		var request NatsRequest
		if err := json.Unmarshal(msg.Data, &request); err != nil {
			handleError(err, msg, &reply)
			return
		}
		var body MapsBodyRequest
		rawBody, err := json.Marshal(request.Body)
		if err != nil {
			handleError(err, msg, &reply)
			return
		}
		if err := json.Unmarshal(rawBody, &body); err != nil {
			handleError(err, msg, &reply)
            return
		}
		coordinates, err := nominatim.GetLatLon(body.Location, body.Limit)

		if err != nil {
			handleError(err, msg, &reply)
			return
		}

		bodyReply := MapsReplyBody{
			Companies:   []MapsCompany{},
			Coordinates: coordinates,
		}

		for _, coordinate := range coordinates {
			companies, err := maps.GetMaps(coordinate, body.Target)
			if err != nil {
				handleError(err, msg, &reply)
				return
			} else {
				for _, company := range companies {
					companyReplay := MapsCompany{
						Title:    company.Title,
						Website:  company.Website,
						Phone:    company.Phone,
						Address:  company.Address,
						PlaceId:  company.PlaceId,
						Category: company.Category,
						Rating:   company.Rating,
					}
					bodyReply.Companies = append(bodyReply.Companies, companyReplay)
				}
			}
		}
		reply.Body = bodyReply
		responseData, err := json.Marshal(reply)
		log.Println("Response data:", string(responseData))
		if err != nil {
			handleError(err, msg, &reply)
			return
		}
		err = msg.Respond(responseData)
		if err != nil {
			handleError(err, msg, &reply)
			return
		}
	})
}

func logSubscribe(subj string, msg *nats.Msg) {
	log.Println("Subscribe:", subj)
	log.Println("Received message:", string(msg.Data))
}

func handleError(err error, msg *nats.Msg, reply *NatsResponse) {
	log.Printf("Error: %v", err)
	reply.Status = 1
	reply.Error = err
	responseData, err := json.Marshal(reply)
	if err != nil {
		log.Printf("Error: Failed to marshal response: %v", err)
		return
	}
	err = msg.Respond(responseData)
	if err != nil {
		log.Printf("Error: Failed to send reply: %v", err)
		return
	}
}