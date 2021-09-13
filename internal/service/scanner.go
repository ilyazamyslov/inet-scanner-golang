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
	Store(string, model.Host) error
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
	if host, ok := s.repo.Load(ip); ok {
		return host, nil
	}
	host, err := scanHost(ip)
	if err != nil {
		return host, err
	}
	err = s.repo.Store(ip, host)
	if err != nil {
		return host, err
	}
	return host, nil
}

func (s *ScannerService) ScanNetwork(ip string, lenMask int) ([]model.Host, error) {
	listHosts, err := listHosts(ip, lenMask)
	if err != nil {
		return []model.Host{{}}, nil
	}
	var result []model.Host
	var listNeedScanHost []string
	for _, val := range listHosts {
		if host, ok := s.repo.Load(val); ok {
			s.logger.Info().Msg("Hit cache " + val)
			result = append(result, host)
		} else {
			listNeedScanHost = append(listNeedScanHost, val)
		}
	}

	scannedHosts, err := scanNetwork(listNeedScanHost)
	if err != nil {
		return result, err
	}

	for _, val := range scannedHosts {
		err := s.repo.Store(val.Ip, val)
		if err != nil {
			return result, err
		}
		result = append(result, val)
	}

	return result, nil
}
