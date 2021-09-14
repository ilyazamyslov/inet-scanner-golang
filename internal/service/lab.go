package service

import (
	"net/http"
	"time"

	"github.com/ilyazamyslov/inet-scanner-golang/internal/model"
	"github.com/rs/zerolog"
)

type LabService struct {
	logger *zerolog.Logger
	repo   Repository
	client HTTPClient
}

func NewLabService(logger *zerolog.Logger, repo Repository) *LabService {
	return &LabService{
		logger: logger,
		repo:   repo,
		client: &http.Client{
			Timeout: time.Duration(time.Minute),
		},
	}
}

func (s *LabService) Query1(start, end string) []model.Host {
	return s.repo.Query1(start, end)
}

func (s *LabService) Query2(start, end string) []model.Host {
	return s.repo.Query2(start, end)
}

func (s *LabService) Query3(start, end string) []model.Host {
	return s.repo.Query3(start, end)
}

func (s *LabService) Query4(start, end string) int {
	return s.repo.Query4(start, end)
}
