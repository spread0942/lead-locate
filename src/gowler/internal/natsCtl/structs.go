package natsCtl

import (
	"time"
)

type NatsRequest struct {
	Method string
	Timestamp time.Time
	Body any
}

type NatsResponse struct {
	Request string
	Method string
	Timestamp time.Time
	Status int
	Body any
	Error error
}

// SUBSCRIBES ****************************************************************

type GowlerBody struct {
	Website string `json:"website"`
}

type GowlerReplyBody struct {
	Site string `json:"site"`
	Domain string `json:"domain"`
	SiteUrls []string `json:"siteUrls"`
	OtherUrls []string `json:"otherUrls"`
	Telephones []string `json:"telephones"`
	Emails []string `json:"emails"`
}