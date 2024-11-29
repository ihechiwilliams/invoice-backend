package appbase

import (
	"github.com/rs/zerolog/log"
	"github.com/samber/do"
	"github.com/samber/lo"
)

type AppBase struct {
	Config      *Config
	ServiceName string
	Injector    *do.Injector
}

func New(options ...func(*AppBase)) *AppBase {
	appBase := &AppBase{}
	for _, o := range options {
		o(appBase)
	}

	return appBase
}

func Init(serviceName string) func(*AppBase) {
	return func(appBase *AppBase) {
		appBase.Config = lo.Must(LoadConfig())
		appBase.ServiceName = serviceName
	}
}

func WithDependencyInjector() func(*AppBase) {
	return func(appBase *AppBase) {
		appBase.Injector = NewInjector(appBase.ServiceName, appBase.Config)
	}
}

func (a *AppBase) Shutdown() {
	err := a.Injector.Shutdown()
	if err != nil {
		log.Panic().Err(err).Msg("injector's shutdown failed")
	}
}
