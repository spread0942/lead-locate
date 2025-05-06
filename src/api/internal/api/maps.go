package api

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/nats-io/nats.go"
)

type MapsInput struct {
	Location string `query:"location" maxLength:"30" example:"Treviso,Veneto,Italia" doc:"Location to search"`
	Target   string `query:"target" maxLength:"30" example:"Pizza" doc:"Target to search"`
	Limit    uint   `query:"limit" example:"10" default:"10" doc:"Number of langitude and longitude to return"`
}

type MapsOutput struct {
	Body struct {
		Companies   []MapsCompany `json:"companies"`
		Coordinates []string      `json:"coordinates"`
	} `json:"body"`
}

type MapsCompany struct {
	Title    string  `json:"title"`
	Website  string  `json:"website"`
	Phone    string  `json:"phone"`
	Address  string  `json:"address"`
	PlaceId  string  `json:"placeId"`
	Category string  `json:"category"`
	Rating   float64 `json:"rating"`
}

type MapsPublish struct {
	Body struct {
		Location string `json:"location"`
		Target   string `json:"target"`
		Limit    uint   `json:"limit"`
	} `json:"body"`
}

func RegisterMaps(api huma.API, nats *nats.Conn) {
	huma.Register(api, huma.Operation{
		OperationID: "maps",
		Summary:     "Maps",
		Method:      http.MethodGet,
		Path:        "/maps",
		Description: "Get maps",
		Tags:        []string{"maps"},
	}, func(ctx context.Context, input *MapsInput) (*MapsOutput, error) {
		resp := &MapsOutput{}
		mapsPublish := MapsPublish{}

		mapsPublish.Body.Location = input.Location
		mapsPublish.Body.Target = input.Target
		mapsPublish.Body.Limit = input.Limit
		jMapsPublish, err := json.Marshal(mapsPublish)

		if err != nil {
			log.Fatal(err)
		}

		reply, err := nats.Request("gowler.request.maps", jMapsPublish, 10*time.Minute)

		if err != nil {
			log.Fatal(err)
		}

		err = json.Unmarshal(reply.Data, &resp)

		if err != nil {
			log.Fatal(err)
		}

		return resp, nil
	})
}
