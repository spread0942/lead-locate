package requests

import (
	"encoding/json"
	"fmt"

	"gowler/internal/gowler"

	"github.com/nats-io/nats.go"
)

type GowlerPublish struct {
	Website string `json:"website"`
}

type GowlerOutput struct {
	Body struct {
		Site string `json:"site"`
		Domain string `json:"domain"`
		SiteUrls []string `json:"siteUrls"`
		OtherUrls []string `json:"otherUrls"`
		Telephones []string `json:"telephones"`
		Emails []string `json:"emails"`
	} `json:"body"`
}

func GowlerIt(msg *nats.Msg) {
	var mp GowlerPublish
	err := json.Unmarshal(msg.Data, &mp)
	if err != nil {
		fmt.Print("Failed to marshal response: %v", err)
	}
	gogo := gowler.NewGowler(mp.Website)
	gogo.Crawl()
	replay := GowlerOutput{}
	replay.Body.Site = gogo.Site
	replay.Body.Domain = gogo.Domain
	replay.Body.SiteUrls = gogo.SiteUrls
	replay.Body.OtherUrls = gogo.OtherUrls
	replay.Body.Telephones = gogo.Telephones
	replay.Body.Emails = gogo.Emails

	responseData, err := json.Marshal(replay)
	if err != nil {
		fmt.Printf("Failed to marshal response: %v", err)
	}
	err = msg.Respond(responseData)
	if err != nil {
		fmt.Printf("Failed to send reply: %v", err)
	}
}