package main

import (
	"github.com/sdmad737/catalog/backend/internal/core/services"
	"github.com/sdmad737/catalog/backend/internal/core/services/reporting/eventbus"
	"github.com/sdmad737/catalog/backend/internal/data/ent"
	"github.com/sdmad737/catalog/backend/internal/data/repo"
	"github.com/sdmad737/catalog/backend/internal/sys/config"
	"github.com/sdmad737/catalog/backend/pkgs/mailer"
)

type app struct {
	conf     *config.Config
	mailer   mailer.Mailer
	db       *ent.Client
	repos    *repo.AllRepos
	services *services.AllServices
	bus      *eventbus.EventBus
}

func new(conf *config.Config) *app {
	s := &app{
		conf: conf,
	}

	s.mailer = mailer.Mailer{
		Host:     s.conf.Mailer.Host,
		Port:     s.conf.Mailer.Port,
		Username: s.conf.Mailer.Username,
		Password: s.conf.Mailer.Password,
		From:     s.conf.Mailer.From,
	}

	return s
}
