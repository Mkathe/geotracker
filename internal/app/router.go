package app

func (s *server) initRoutes() {
	s.app.Get("/healthz", s.CheckHealth)
}
