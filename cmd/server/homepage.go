package server

import (
	"net/http"
)

type homepageResponse struct {
	SwaggerUI string
}

func (s *Server) homepage(w http.ResponseWriter, r *http.Request) {
	response := homepageResponse{
		SwaggerUI: "http://127.0.0.1:8082/?url=http://" + r.Host + "/swagger.yaml",
	}
	s.io.writeJSON(w, http.StatusOK, response)
}
