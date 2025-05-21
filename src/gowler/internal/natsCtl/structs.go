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

type MapsBodyRequest struct {
	Location string `json:"location"`
	Target   string `json:"target"`
	Limit    uint   `json:"limit"`
}

type MapsReplyBody struct {
	Companies   []MapsCompany `json:"companies"`
	Coordinates []string      `json:"coordinates"`
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