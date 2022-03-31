package server

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/ReLium/crud/internal/repository"
	"github.com/go-chi/chi"
)

type Server struct {
	repository repository.Repository
	io         *IO
}

type cat struct {
	Name       string `json:"name"`
	Gender     string `json:"gender"`
	Color      string `json:"color"`
	Vaccinated bool   `json:"vaccinated"`
}

type responseList struct {
	Items []cat `json:"items"`
}

func NewServer(repository repository.Repository) *Server {
	return &Server{
		repository: repository,
		io:         NewIO(),
	}
}

func (s *Server) Serve(addr string) error {
	r := chi.NewRouter()

	// @TODO: Setup CORS settings properly. Now enabled for all for testing demo purposes.
	r.Use(s.allowCORS)

	// Log each request.
	r.Use(s.requestLogger)
	r.Get("/", s.homepage)
	r.Get("/swagger.yaml", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		http.ServeFile(w, r, "./doc/swagger.yaml")
	})
	r.Route("/cats", func(r chi.Router) {
		r.Get("/", s.listCats)
		r.Route("/{CatName}", func(r chi.Router) {
			r.Get("/", s.getCat)
			r.Put("/", s.updateCat)
			r.Delete("/", s.deleteCat)
		})
		r.Route("/add", func(r chi.Router) {
			r.Post("/", s.insertCat)
		})
	})
	return http.ListenAndServe(addr, r)
}

func (s *Server) listCats(w http.ResponseWriter, r *http.Request) {
	catParams := s.parseCatParams(r)
	query := repository.Query(*catParams)

	cats, err := s.repository.Find(&query)
	if err != nil {
		err = s.io.writeError(w, err)
		if err != nil {
			log.Println(err)
		}
		return
	}
	response := responseList{
		Items: make([]cat, 0, len(cats)),
	}
	for _, c := range cats {
		response.Items = append(response.Items, cat(*c))
	}
	err = s.io.writeJSON(w, http.StatusOK, response)
	if err != nil {
		log.Println(err)
	}
}

func (s *Server) getCat(w http.ResponseWriter, r *http.Request) {
	cat, err := s.repository.Get(chi.URLParam(r, "CatName"))
	if err != nil {
		s.io.writeError(w, err)
		return
	}
	s.io.writeJSON(w, http.StatusOK, cat)
}

func (s *Server) updateCat(w http.ResponseWriter, r *http.Request) {
	catParams := s.parseCatParams(r)
	catUpdate := &repository.CatUpdate{
		Name:       chi.URLParam(r, "CatName"),
		Gender:     catParams.Gender,
		Color:      catParams.Color,
		Vaccinated: catParams.Vaccinated,
	}
	err := s.repository.Update(catUpdate)
	if err != nil {
		s.io.writeError(w, err)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (s *Server) deleteCat(w http.ResponseWriter, r *http.Request) {
	err := s.repository.Delete(chi.URLParam(r, "CatName"))
	if err != nil {
		s.io.writeError(w, err)
		return
	}
}

func (s *Server) insertCat(w http.ResponseWriter, r *http.Request) {
	var cat cat
	err := s.io.readJSON(w, r, &cat)
	if err != nil {
		s.io.writeError(w, err)
		return
	}
	repoCat := repository.Cat(cat)

	err = s.repository.Insert(&repoCat)
	if err != nil {
		s.io.writeError(w, err)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (s *Server) requestLogger(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		defer func() {
			log.Printf("Incoming HTTP request: %s %s in %s", r.Method, r.URL, time.Since(start))
		}()

		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

type catParams struct {
	Gender     string
	Color      string
	Vaccinated *bool
}

func (s *Server) parseCatParams(r *http.Request) *catParams {
	catParams := &catParams{
		Gender:     r.URL.Query().Get("gender"),
		Color:      r.URL.Query().Get("color"),
		Vaccinated: nil,
	}
	if r.URL.Query().Has("vaccinated") {
		vaccinated, err := strconv.ParseBool(r.URL.Query().Get("vaccinated"))
		if err == nil {
			catParams.Vaccinated = &vaccinated
		}
	}
	return catParams
}

func (s *Server) allowCORS(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "*")
		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
