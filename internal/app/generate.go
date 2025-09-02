package app

import (
	"github.com/gofiber/fiber/v2"
	"github.com/hashicorp/go-hclog"
	"github.com/magzhan/geotracker/pkg/config"
	"github.com/magzhan/geotracker/pkg/db/postgres"
	"github.com/magzhan/geotracker/pkg/db/redis"
)

func (s *server) generate() error {
	s.app = fiber.New(
		fiber.Config{
			DisableStartupMessage: true,
		})

	s.logger = hclog.New(&hclog.LoggerOptions{
		Name:               "library-api",
		JSONFormat:         true,
		JSONEscapeDisabled: true,
		Level:              hclog.Debug,
	})

	db, err := postgres.LoadDatabase(config.Get().ConnStr)
	if err != nil {
		s.logger.Error("Error loading database", "error", err)
		return err
	}
	s.db = db
	s.cache = redis.LoadCache("", "", 0) // Fix it

	//kcAuth, err := auth.StartKeyCloakAuth(config.Get().KeycloakURL, config.Get().KeycloakRole)
	//if err != nil {
	//	s.logger.Error("Error starting keycloak auth", "error", err)
	//	return err
	//}
	//
	//s.useMiddleware(kcAuth)

	s.initRoutes()

	return nil
}
