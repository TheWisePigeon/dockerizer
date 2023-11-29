package main

import (
	"dockerizer/server"
	"net/http"
	"github.com/docker/docker/client"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

type Payload struct {
	Action string `json:"action"`
	Entity string `json:"entity"`
}

func main() {
	dockerClient, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}
	defer dockerClient.Close()
	r := chi.NewRouter()
	validate := validator.New(validator.WithRequiredStructEnabled())
  server := server.NewServer(dockerClient, validate)

  server.RegisterRoutes(r)

	err = http.ListenAndServe(":6061", r)
	if err != nil {
		panic(err)
	}
}
