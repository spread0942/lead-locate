package natsUtils

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