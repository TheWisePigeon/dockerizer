package tests

import (
	"dockerizer/server"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/docker/docker/client"
	"github.com/go-playground/validator/v10"
)

func TestImages(t *testing.T) {
	client, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	validate := validator.New(validator.WithRequiredStructEnabled())
	if err != nil {
		t.Errorf("Error while creating docker client %q", err)
	}
	server := server.NewServer(client, validate)
	t.Run("Get list of all images", func(t *testing.T) {
		req, err := http.NewRequest(
			"GET",
			"/image",
			nil,
		)
		if err != nil {
			t.Errorf("Failed to create request %q", err)
		}
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.GetAllImages)
		handler.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("Wanted %q got %q", http.StatusOK, rr.Code)
		}
	})

}
