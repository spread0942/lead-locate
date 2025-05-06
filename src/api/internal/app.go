package internal

import (
	"context"
	"os"
	"os/signal"

	"github.com/nats-io/nats.go"
)

type App struct {
	SrvCtx context.Context
	Stop   context.CancelFunc
}

type appContextKey string

func NewApp(name string, version string) *App {
	srvctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	
	srvctx = context.WithValue(srvctx, appContextKey("serviceName"), name)
	srvctx = context.WithValue(srvctx, appContextKey("serviceVersion"), version)
	return &App{
		SrvCtx: srvctx,
		Stop:   stop,
	}
}

func (a *App) SetNatsContext(natsUrl string) error {
	nats, err := nats.Connect(natsUrl)
	a.SrvCtx = context.WithValue(a.SrvCtx, appContextKey("nats"), nats)
	return err
}

func (a *App) GetNatsConnection() *nats.Conn {
	return a.SrvCtx.Value(appContextKey("nats")).(*nats.Conn)
}

func (a *App) StartApiServer(apiUrl string) {
	go StartApiServer(*a, apiUrl)
}