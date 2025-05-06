package requests

import (
	"encoding/json"
	"fmt"
	"log"

	"gowler/internal/google/maps"
	"gowler/internal/nominatim"

	"github.com/nats-io/nats.go"
)

type MapsPublish struct {
	Body struct {
		Location string `json:"location" example:"Treviso,Veneto,Italia" doc:"Location to search"`
		Target   string `json:"target" example:"Pizza" doc:"Target to search"`
		Limit    uint   `json:"limit" example:"10" doc:"Number of langitude and longitude to return" default:"10"`
	} `json:"body"`
}

type MapsResponse struct {
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

func GetMaps(msg *nats.Msg) {
	var mp MapsPublish
	err := json.Unmarshal(msg.Data, &mp)
	if err != nil {
		log.Fatal(err)
	}
	if err != nil {
		log.Fatal(err)
	}
	coordinates, err := nominatim.GetLatLon(mp.Body.Location, mp.Body.Limit)

	if err != nil {
		fmt.Println(err)
	}

	replay := MapsResponse{}
	replay.Body.Coordinates = coordinates

	for _, coordinate := range coordinates {
		companies, err := maps.GetMaps(coordinate, mp.Body.Target)
		if err != nil {
			fmt.Println(err)
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
				replay.Body.Companies = append(replay.Body.Companies, companyReplay)
			}
		}
	}

	responseData, err := json.Marshal(replay)
	if err != nil {
		fmt.Printf("Failed to marshal response: %v", err)
	}
	err = msg.Respond(responseData)
	if err != nil {
		fmt.Printf("Failed to send reply: %v", err)
	}
}
