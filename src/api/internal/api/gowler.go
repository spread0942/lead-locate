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
	"api/internal/utils"
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

type GowlerBodyRequest struct {
	Website string `json:"website"`
}

func RegisterGowler(api huma.API, nats *utils.NatsCtl) {
	huma.Register(api, huma.Operation{
		OperationID: "gowler",
		Summary:     "Crawler a website",
		Method:      http.MethodGet,
		Path:        "/gowler",
	}, func(ctx context.Context, input *GowlerInput) (*GowlerOutput, error) {
		resp := &GowlerOutput{}
		body := GowlerBodyRequest{
			Website: input.Website,
		}
		data, err := nats.Request("gowler.request.gowler", "gowler", body)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal response: %w", err)
		}
		jsonData, err := json.Marshal(data)
		err = json.Unmarshal(jsonData, &resp)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal to GowlerOutput: %w", err)
		}
		return resp, nil
	})
}
