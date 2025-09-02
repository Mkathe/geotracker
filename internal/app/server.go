package app

import (
	"database/sql"
	"github.com/gofiber/fiber/v2"
	"github.com/hashicorp/go-hclog"
	"github.com/magzhan/geotracker/pkg/config"
	"github.com/redis/go-redis/v9"
)

type server struct {
	app    *fiber.App
	logger hclog.Logger
	db     *sql.DB
	cache  *redis.Client
}

func Run() error {
	var err error

	s := new(server)

	err = s.generate()
	if err != nil {
		return err
	}

	errChan := make(chan error, 1)
	go func() {
		err = s.app.Listen(config.Get().Ip + ":" + config.Get().Port)
		if err != nil {
			errChan <- err
			s.logger.Error("Error starting server", "error", err)
		}
	}()
	s.gracefulShutdown(s.logger)

	select {
	case err = <-errChan:
		return err
	}
}
