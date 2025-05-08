package internal

import (
	"context"
	"os"
	"os/signal"

	"api/internal/utils"
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
	natsCtl, err := utils.NewNatsCtl(natsUrl)
	a.SrvCtx = context.WithValue(a.SrvCtx, appContextKey("nats"), natsCtl)
	return err
}

func (a *App) GetNatsConnection() *utils.NatsCtl {
	return a.SrvCtx.Value(appContextKey("nats")).(*utils.NatsCtl)
}

func (a *App) StartApiServer(apiUrl string) {
	go StartApiServer(*a, apiUrl)
}