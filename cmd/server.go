package cmd

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/ReLium/crud/internal/mongodb"
	"github.com/ReLium/crud/internal/repository"
	"github.com/go-chi/chi"
)

const DefaultMongoUrl = "mongodb://admin:pass@127.0.0.1:27017/"
const DefaultMongoTimeoutMilliseconds = 1000

type cat struct {
	Name       string `json:"name"`
	Gender     string `json:"gender"`
	Color      string `json:"color"`
	Vaccinated bool   `json:"vaccinated"`
}

type responseList struct {
	Items []cat `json:"items"`
}

type serverContext struct {
	repository repository.Repository
	io         *IO
}

// Execute executes the root command.
func Server() error {

	mongoDBClient, err := mongodb.NewClient(DefaultMongoUrl, DefaultMongoTimeoutMilliseconds)
	if err != nil {
		return err
	}
	repo := repository.NewMongoDBRepo(mongoDBClient)
	serverContext := &serverContext{
		repository: repo,
		io:         NewIO(),
	}

	r := chi.NewRouter()

	// Log each request.
	r.Use(serverContext.requestLogger)

	r.Route("/cats", func(r chi.Router) {
		r.Get("/", serverContext.listCats)
		r.Route("/{CatName}", func(r chi.Router) {
			r.Get("/", serverContext.getCat)
			r.Put("/", serverContext.updateCat)
			r.Delete("/", serverContext.deleteCat)
		})
		r.Route("/add", func(r chi.Router) {
			r.Post("/", serverContext.insertCat)
		})
	})
	return http.ListenAndServe(":8080", r)
}

func (s *serverContext) listCats(w http.ResponseWriter, r *http.Request) {
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

func (s *serverContext) getCat(w http.ResponseWriter, r *http.Request) {
	cat, err := s.repository.Get(chi.URLParam(r, "CatName"))
	if err != nil {
		s.io.writeError(w, err)
		return
	}
	s.io.writeJSON(w, http.StatusOK, cat)
}

func (s *serverContext) updateCat(w http.ResponseWriter, r *http.Request) {
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

func (s *serverContext) deleteCat(w http.ResponseWriter, r *http.Request) {
	err := s.repository.Delete(chi.URLParam(r, "CatName"))
	if err != nil {
		s.io.writeError(w, err)
		return
	}
}

func (s *serverContext) insertCat(w http.ResponseWriter, r *http.Request) {
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

func (s *serverContext) requestLogger(h http.Handler) http.Handler {
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

func (s *serverContext) parseCatParams(r *http.Request) *catParams {
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
