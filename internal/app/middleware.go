package app

import (
	"github.com/ansrivas/fiberprometheus/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/magzhan/geotracker/pkg/auth"
	"github.com/magzhan/geotracker/pkg/ws"
)

func (s *server) useMiddleware(kcAuth *auth.KeyCloakAuth) {
	s.app.Use(cors.New(cors.Config{
		AllowMethods:     "GET, POST, PATCH, DELETE",
		AllowHeaders:     "Origin, Content-Type, Accept",
		AllowCredentials: false,
		AllowOrigins:     "*",
	}))

	prometheus := fiberprometheus.New("library-service")
	prometheus.RegisterAt(s.app, "/metrics")
	s.app.Use(prometheus.Middleware)

	s.app.Use(kcAuth.Auth)

	s.app.Use(ws.WebsocketsCheckMiddleware)
}
