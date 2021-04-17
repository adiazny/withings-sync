package api

func (s *Server) routes() {
	s.router.HandleFunc("/about", s.handleAboutEndpoint())
	s.router.HandleFunc("/weight", s.handleWithingsCallback())
}
