package service

import (
	"net/http"
	"time"

	"github.com/ilyazamyslov/inet-scanner-golang/internal/model"
	"github.com/rs/zerolog"
)

type ScannerService struct {
	logger *zerolog.Logger
	repo   Repository
	client HTTPClient
}

type Repository interface {
	Load(string) (model.Host, bool)
	Store(string, model.Host)
}

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

func NewScannerService(logger *zerolog.Logger, repo Repository) *ScannerService {
	return &ScannerService{
		logger: logger,
		repo:   repo,
		client: &http.Client{
			Timeout: time.Duration(time.Minute),
		},
	}
}

func (s *ScannerService) ScanHost(ip string) (model.Host, error) {
	return scanHost(ip)
}

func (s *ScannerService) ScanNetwork(ip string, lenMask int) ([]model.Host, error) {
	/*req, err := http.NewRequest(http.MethodGet, fmt.Sprintf(stockAPIFormat, ticker, s.apiKey), nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()*/

	return scanNetwork(ip, lenMask)
}
