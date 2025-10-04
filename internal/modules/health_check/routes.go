package healthCheck

func (s *healthCheckService) SetupRoutes() {
	healthCheck := s.route.Group("/health-check")
	healthCheck.Get("/", s.Check)
}
