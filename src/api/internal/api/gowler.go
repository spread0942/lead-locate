package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"api/internal/natsCtl"
)

type GowlerInput struct {
	Website string `query:"website" maxLength:"30" example:"https://www.google.com" doc:"Website to crawl" required:"true"`
}

type GowlerReplyBody struct {
	Site string `json:"site"`
	Domain string `json:"domain"`
	SiteUrls []string `json:"siteUrls"`
	OtherUrls []string `json:"otherUrls"`
	Telephones []string `json:"telephones"`
	Emails []string `json:"emails"`
}

type GowlerBodyRequest struct {
	Website string `json:"website"`
}

type ApiResponse struct {
	Body struct {
		Site string `json:"site"`
		Domain string `json:"domain"`
		SiteUrls []string `json:"siteUrls"`
		OtherUrls []string `json:"otherUrls"`
		Telephones []string `json:"telephones"`
		Emails []string `json:"emails"`
	}
}

func RegisterGowler(api huma.API, nats *natsCtl.NatsCtl) {
	huma.Register(api, huma.Operation{
		OperationID: "gowler",
		Summary:     "Crawler a website",
		Method:      http.MethodGet,
		Path:        "/gowler",
		Description: "Crawl a website and return the results",
		Tags:        []string{"Gowler"},
	}, func(ctx context.Context, input *GowlerInput) (*ApiResponse, error) {
		if input.Website == "" {
			return nil, huma.NewError(http.StatusBadRequest, "missing website parameter")
		}
		body := GowlerBodyRequest{
			Website: input.Website,
		}
		data, err := nats.Request("gowler.request.gowler", "gowler", body)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal response: %w", err)
		}
		reply := &GowlerReplyBody{}
		err = json.Unmarshal(data, &reply)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal to GowlerOutput: %w", err)
		}
		resp := &ApiResponse{}
		resp.Body.Site = reply.Site
		resp.Body.Domain = reply.Domain
		resp.Body.SiteUrls = reply.SiteUrls
		resp.Body.OtherUrls = reply.OtherUrls
		resp.Body.Telephones = reply.Telephones
		resp.Body.Emails = reply.Emails
		return resp, nil
	})
}
