package tests

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestImages(t *testing.T) {
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

	})

}
