package app

import "github.com/gofiber/contrib/websocket"

func (s *server) initRoutes() {
	s.app.Get("/healthz", s.CheckHealth)
	s.app.Get("/location", websocket.New(s.ServeWs))
}
