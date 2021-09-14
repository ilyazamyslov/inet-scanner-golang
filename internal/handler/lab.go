package handler

import (
	"net/http"

	"github.com/ilyazamyslov/inet-scanner-golang/internal/model"
	"github.com/rs/zerolog"
)

const (
	LabPath = "/lab"
)

type Lab struct {
	logger  *zerolog.Logger
	service LabService
}

type LabService interface {
	Query1(start, end string) []model.Host
	Query2(start, end string) []model.Host
	Query3(start, end string) []model.Host
	Query4(start, end string) int
}

func NewLab(logger *zerolog.Logger, service LabService) *Lab {
	return &Lab{
		logger:  logger,
		service: service,
	}
}

func (h *Lab) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var result model.LabResp
	result.Q1 = h.service.Query1("5.255.0.0", "5.255.255.25")
	result.Q2 = h.service.Query2("5.255.0.0", "5.255.255.25")
	result.Q3 = h.service.Query3("5.255.0.0", "5.255.255.25")
	result.Q4 = h.service.Query4("5.255.0.0", "5.255.255.25")

	writeResponse(w, http.StatusOK, result)
}
