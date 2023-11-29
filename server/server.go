package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

type Server struct {
	dockerClient *client.Client
	validate     *validator.Validate
}

func NewServer(client *client.Client, validate *validator.Validate) *Server {
	return &Server{
		dockerClient: client,
		validate:     validate,
	}
}

func (s *Server) GetAllImages(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	images, err := s.dockerClient.ImageList(ctx, types.ImageListOptions{})
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	data, err := json.Marshal(images)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(data)
	return
}

func (s *Server) RegisterRoutes(r chi.Router) {
	r.Route("/image", func(r chi.Router) {
		r.Get("/", s.GetAllImages)
	})
}
