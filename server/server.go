package server

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
)

type Server struct {
	dockerClient *client.Client
	validate     *validator.Validate
	logger       *logrus.Logger
}

func NewServer(client *client.Client, validate *validator.Validate, logger *logrus.Logger) *Server {
	return &Server{
		dockerClient: client,
		validate:     validate,
		logger:       logger,
	}
}

func (s *Server) GetAllImages(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	images, err := s.dockerClient.ImageList(ctx, types.ImageListOptions{})
	if err != nil {
		s.logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	data, err := json.Marshal(images)
	if err != nil {
		s.logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(data)
	return
}

func (s *Server) GetImageByName(w http.ResponseWriter, r *http.Request) {
	imageName := chi.URLParam(r, "image")
	//For tests only
	imageName = strings.Split(r.URL.String(), "/")[2]
	if imageName == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var image types.ImageSummary
	found := false
	ctx := context.Background()
	images, err := s.dockerClient.ImageList(ctx, types.ImageListOptions{})
	if err != nil {
		s.logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	for _, img := range images {
		for _, tag := range img.RepoTags {
			if tag == imageName {
				image = img
				found = true
				break
			}
		}
	}
	if found {
		data, err := json.Marshal(image)
		if err != nil {
			s.logger.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(data)
		return
	}
	w.WriteHeader(http.StatusNotFound)
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
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(errStr))
		return
	}
	defer out.Close()
	var buffer bytes.Buffer
	io.Copy(&buffer, out)
	w.WriteHeader(http.StatusOK)
	return
}

func (s *Server) RemoveImage(w http.ResponseWriter, r *http.Request) {
	image := chi.URLParam(r, "image")
	//For tests only
	image = strings.Split(r.URL.String(), "/")[3]
	if image == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	ctx := context.Background()
	_, err := s.dockerClient.ImageRemove(ctx, image, types.ImageRemoveOptions{})
	if err != nil {
    s.logger.Error(err)
		errStr := err.Error()
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(errStr))
		return
	}
	w.WriteHeader(http.StatusOK)
	return
}

func (s *Server) GetAllContainers(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	containers, err := s.dockerClient.ContainerList(ctx, types.ContainerListOptions{})
	if err != nil {
		s.logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	data, err := json.Marshal(containers)
	w.WriteHeader(http.StatusOK)
	w.Write(data)
	return
}

func (s *Server) CreateContainer() {

}

func (s *Server) RegisterRoutes(r chi.Router) {
	r.Route("/image", func(r chi.Router) {
		r.Get("/", s.GetAllImages)
		r.Get("/{image}", s.GetImageByName)
		r.Get("/pull/{image}", s.PullImage)
		r.Delete("/remove/{image}", s.RemoveImage)
	})
}
