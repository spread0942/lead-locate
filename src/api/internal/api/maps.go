package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"api/internal/natsCtl"
)

type MapsCompany struct {
	Title    string  `json:"title"`
	Website  string  `json:"website"`
	Phone    string  `json:"phone"`
	Address  string  `json:"address"`
	PlaceId  string  `json:"placeId"`
	Category string  `json:"category"`
	Rating   float64 `json:"rating"`
}

// api structs
type MapsInput struct {
	Location string `query:"location" maxLength:"30" example:"Treviso,Veneto,Italia" doc:"Location to search" required:"true"`
	Target   string `query:"target" maxLength:"30" example:"Pizza" doc:"Target to search" required:"true"`
	Limit    uint   `query:"limit" example:"10" default:"10" doc:"Number of langitude and longitude to return"`
}

type ApiOutput struct {
	Body struct {
		Companies   []MapsCompany `json:"companies"`
		Coordinates []string      `json:"coordinates"`
	} `json:"body"`
}

// nats structs
type MapsReplyBody struct {
	Companies   []MapsCompany `json:"companies"`
	Coordinates []string      `json:"coordinates"`
}

type MapsBodyRequest struct {
	Location string `json:"location"`
	Target   string `json:"target"`
	Limit    uint   `json:"limit"`
}

func RegisterMaps(api huma.API, nats *natsCtl.NatsCtl) {
	huma.Register(api, huma.Operation{
		OperationID: "maps",
		Summary:     "Maps",
		Method:      http.MethodGet,
		Path:        "/maps",
		Description: "Get maps",
		Tags:        []string{"Maps"},
	}, func(ctx context.Context, input *MapsInput) (*ApiOutput, error) {
		if input.Location == "" {
			return nil, huma.NewError(http.StatusBadRequest, "missing location parameter")
		}
		if input.Target == "" {
			return nil, huma.NewError(http.StatusBadRequest, "missing target parameter")
		}
		if input.Limit < 1 {
			return nil, huma.NewError(http.StatusBadRequest, "limit must be greater than 0")
		}
		if input.Limit > 10 {
			return nil, huma.NewError(http.StatusBadRequest, "limit must be less than 10")
		}
		body := MapsBodyRequest{
			Location: input.Location,
			Target:   input.Target,
			Limit:    input.Limit,
		}
		data, err := nats.Request("gowler.request.maps", "maps", body)
		fmt.Printf("data: %s\n", string(data))
		
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal response: %w", err)
		}
		reply := &MapsReplyBody{}
		err = json.Unmarshal(data, &reply)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal to GowlerOutput: %w", err)
		}
		resp := &ApiOutput{}
		resp.Body.Companies = reply.Companies
		resp.Body.Coordinates = reply.Coordinates

		return resp, nil
	})
}
