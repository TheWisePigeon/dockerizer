package tests

import (
	"dockerizer/server"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"github.com/docker/docker/api/types"
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
	imageId := ""
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
			t.Errorf("Wanted %v got %v", http.StatusOK, rr.Code)
		}
	})

	t.Run("Pull image", func(t *testing.T) {
		imageName := "alpine"
		req, err := http.NewRequest(
			"GET",
			fmt.Sprintf("/image/pull/%s", imageName),
			nil,
		)
		if err != nil {
			t.Errorf("Failed to create request %q", err)
		}
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.PullImage)
		handler.ServeHTTP(rr, req)
		if rr.Code != http.StatusOK {
			t.Errorf("Wanted %v got %v", http.StatusOK, rr.Code)
		}
	})

	t.Run("Pull non existing image", func(t *testing.T) {
		imageName := "a-non-existing-image"
		req, err := http.NewRequest(
			"GET",
			fmt.Sprintf("/image/pull/%s", imageName),
			nil,
		)
		if err != nil {
			t.Errorf("Failed to create request %q", err)
		}
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.PullImage)
		handler.ServeHTTP(rr, req)
		if rr.Code != http.StatusBadRequest {
			t.Errorf("Wanted %v got %v", http.StatusBadRequest, rr.Code)
		}
	})

	t.Run("Get image ny name", func(t *testing.T) {
		imageName := "alpine:latest"
		req, err := http.NewRequest(
			"GET",
			fmt.Sprintf("/image/%s", imageName),
			nil,
		)
		if err != nil {
			t.Errorf("Failed to create request %q", err)
		}
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.GetImageByName)
		handler.ServeHTTP(rr, req)
		if rr.Code != http.StatusOK {
			t.Errorf("Wanted %v got %v", http.StatusOK, rr.Code)
		}
		respBytes, err := io.ReadAll(rr.Body)
		if err != nil {
			t.Errorf("Failed to read response body %q", err)
		}
		var image types.ImageSummary
		if err := json.Unmarshal(respBytes, &image); err != nil {
			t.Errorf("Failed to unmarshal response %q", err)
		}
		imageId = image.ID
		if imageId == "" {
			t.Errorf("Wanted types.ImageSummary got %v", image)
		}
	})

	t.Run("Get non existing image", func(t *testing.T) {
		imageName := "a-non-existing-image"
		req, err := http.NewRequest(
			"GET",
			fmt.Sprintf("/image/%s", imageName),
			nil,
		)
		if err != nil {
			t.Errorf("Failed to create request %q", err)
		}
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.GetImageByName)
		handler.ServeHTTP(rr, req)
		if rr.Code != http.StatusNotFound {
			t.Errorf("Wanted %v got %v", http.StatusNotFound, rr.Code)
		}
	})

	t.Run("remove image", func(t *testing.T) {
		req, err := http.NewRequest(
			"DELETE",
			fmt.Sprintf("/image/remove/%s", imageId),
			nil,
		)
		if err != nil {
			t.Errorf("Failed to create request %q", err)
		}
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.RemoveImage)
		handler.ServeHTTP(rr, req)
		if rr.Code != http.StatusOK {
			t.Errorf("Wanted %v got %v with error %s", http.StatusOK, rr.Code, rr.Body.String())
		}
	})

	t.Run("delete non existing image", func(t *testing.T) {
		req, err := http.NewRequest(
			"DELETE",
			fmt.Sprintf("/image/remove/%s", "a-non-existing-image"),
			nil,
		)
		if err != nil {
			t.Errorf("Failed to create request %q", err)
		}
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.RemoveImage)
		handler.ServeHTTP(rr, req)
		if rr.Code != http.StatusBadRequest {
			t.Errorf("Wanted %v got %v ", http.StatusBadRequest, rr.Code)
		}
	})
}
