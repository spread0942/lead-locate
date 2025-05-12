package internal

import (
	"context"
	"os"
	"os/signal"

	"gowler/internal/utils"
)

type App struct {
	SrvCtx context.Context
	Stop   context.CancelFunc
}

type appContextKey string

// configurate the application
func NewApp(name string, version string) *App {
	srvctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)

	srvctx = context.WithValue(srvctx, appContextKey("serviceName"), name)
	srvctx = context.WithValue(srvctx, appContextKey("serviceVersion"), version)
	return &App{
		SrvCtx: srvctx,
		Stop:   stop,
	}
}

// set nats in the context
func (a *App) SetNatsContext(natsUrl string) error {
	natsCtl, err := utils.NewNatsCtl(natsUrl)
	a.SrvCtx = context.WithValue(a.SrvCtx, appContextKey("nats"), natsCtl)
	return err
}

// an handler to get back natsCtl
func (a *App) GetNatsConnection() *utils.NatsCtl {
	return a.SrvCtx.Value(appContextKey("nats")).(*utils.NatsCtl)
}

// start the nats subscribes
func (a *App) StartNatsSubscibes() {
	serverName := a.SrvCtx.Value(appContextKey("serviceName")).(string)
	natsCtl := a.GetNatsConnection()
	natsCtl.SubscribeHandler(serverName)
}
