package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ReLium/crud/internal/repository"
)

func TestServerHandlers_UpdateCat(t *testing.T) {
	bTrue := true
	tests := []struct {
		name      string
		url       string
		catUpdate *repository.CatUpdate
	}{
		{
			name:      "No updates",
			url:       "/cats/Felix",
			catUpdate: &repository.CatUpdate{},
		},
		{
			name: "Change color",
			url:  "/cats/Felix?color=red",
			catUpdate: &repository.CatUpdate{
				Color: "red",
			},
		},
		{
			name: "Change color and vaccination",
			url:  "/cats/Felix?color=red&vaccinated=true",
			catUpdate: &repository.CatUpdate{
				Color:      "red",
				Vaccinated: &bTrue,
			},
		},
		{
			name: "Invalid vaccinated",
			url:  "/cats/Felix?color=red&vaccinated=tru",
			catUpdate: &repository.CatUpdate{
				Color: "red",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("PUT", tt.url, nil)
			if err != nil {
				t.Fatal(err)
			}
			rr := httptest.NewRecorder()

			mockRepo := new(repository.MockRepository)
			mockRepo.On("Update", tt.catUpdate).Return(nil).Times(1)

			server := NewServer(mockRepo)
			http.HandlerFunc(server.updateCat).ServeHTTP(rr, req)
			mockRepo.AssertExpectations(t)
		})
	}
}
