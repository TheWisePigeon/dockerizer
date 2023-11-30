package main

import (
	"dockerizer/server"
	"fmt"
	"net/http"
	"os"

	"github.com/docker/docker/client"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
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
	logger := logrus.New()
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	logfile := fmt.Sprintf("%s/.dockerizer.log", userHomeDir)
	f, err := os.OpenFile(logfile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
  defer f.Close()
	if err != nil {
		panic(err)
	}
	logger.SetReportCaller(true)
  logger.SetOutput(f)
	server := server.NewServer(dockerClient, validate, logger)

	server.RegisterRoutes(r)

	err = http.ListenAndServe(":6061", r)
	if err != nil {
		panic(err)
	}
}
