package server

func (s HttpServer) addRoutes() {
	s.router.HandleFunc("POST /register", s.authHandler.registerHandler)
	s.router.HandleFunc("POST /refresh", s.authHandler.refreshHandler)
}
