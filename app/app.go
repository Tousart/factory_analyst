package app

import (
	ranalyst "github.com/tousart/factory_analyst/internal/repository/analyst"
	sanalyst "github.com/tousart/factory_analyst/internal/service/analyst"
	tanalyst "github.com/tousart/factory_analyst/internal/transport/analyst"
	"github.com/tousart/factory_analyst/pkg/server"
)

type service struct {
	analystSrv *sanalyst.Service
}

type repository struct {
	analystRepo *ranalyst.Repository
}

type App struct {
	// Transports все gRPC-транспорты, регистрируемые на сервере.
	Transports []server.Service

	// service все сервисные слои.
	service service

	// repository все слои репозитория.
	repository repository
}

func NewApp() *App {
	app := &App{}

	app.initRepository()
	app.initService()
	app.initTransport()

	return app
}

func (app *App) initTransport() {
	app.Transports = []server.Service{
		tanalyst.NewTransport(app.service.analystSrv),
	}
}

func (app *App) initService() {
	app.service = service{
		analystSrv: sanalyst.NewService(app.repository.analystRepo),
	}
}

func (app *App) initRepository() {
	analystRepo := ranalyst.NewRepository()

	app.repository = repository{
		analystRepo: analystRepo,
	}
}
