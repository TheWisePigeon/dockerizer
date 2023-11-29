package server

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

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

func (s *Server) PullImage(w http.ResponseWriter, r *http.Request) {
	image := chi.URLParam(r, "image")
  //For tests only
  image = strings.Split(r.URL.String(), "/")[3]
	if image == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	ctx := context.Background()
	out, err := s.dockerClient.ImagePull(ctx, image, types.ImagePullOptions{})
	if err != nil {
    errStr := err.Error()
	if strings.Contains(errStr, "requested access to the resource is denied") {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer out.Close()
	var buffer bytes.Buffer
	io.Copy(&buffer, out)
	w.WriteHeader(http.StatusOK)
	return
}

func (s *Server) RegisterRoutes(r chi.Router) {
	r.Route("/image", func(r chi.Router) {
		r.Get("/", s.GetAllImages)
		r.Get("/pull/{image}", s.PullImage)
	})
}
