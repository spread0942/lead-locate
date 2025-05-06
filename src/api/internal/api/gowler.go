package api

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/nats-io/nats.go"
)

type GowlerInput struct {
	Website string `query:"website" maxLength:"30" example:"https://www.google.com" doc:"Website to crawl"`
}

type GowlerOutput struct {
	Body struct {
		Site       string   `json:"site"`
		Domain     string   `json:"domain"`
		SiteUrls   []string `json:"siteUrls"`
		OtherUrls  []string `json:"otherUrls"`
		Telephones []string `json:"telephones"`
		Emails     []string `json:"emails"`
	} `json:"body"`
}

type GowlerPublish struct {
	Website string `json:"website"`
}

func RegisterGowler(api huma.API, nats *nats.Conn) {
	huma.Register(api, huma.Operation{
		OperationID: "gowler",
		Summary:     "Crawler a website",
		Method:      http.MethodGet,
		Path:        "/gowler",
	}, func(ctx context.Context, input *GowlerInput) (*GowlerOutput, error) {
		resp := &GowlerOutput{}
		mapsPublish := GowlerPublish{
			Website: input.Website,
		}
		jMapsPublish, err := json.Marshal(mapsPublish)
		if err != nil {
			fmt.Println(err)
		}
		reply, err := nats.Request("gowler.request.gowler", jMapsPublish, 10*time.Minute)
		if err != nil {
			log.Fatal(err)
		}
		err = json.Unmarshal(reply.Data, resp)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal response: %w", err)
		}
		return resp, nil
	})
}
