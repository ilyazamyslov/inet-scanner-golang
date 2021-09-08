package handler

import (
	"fmt"
	"net"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/ilyazamyslov/inet-scanner-golang/internal/model"
	"github.com/rs/zerolog"
)

const (
	ScanHostPath = "/host/{ip}"
)

type HostScan struct {
	logger  *zerolog.Logger
	service ScannerService
}

func NewHostScan(logger *zerolog.Logger, service ScannerService) *HostScan {
	return &HostScan{
		logger:  logger,
		service: service,
	}
}

func (h *HostScan) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ip := chi.URLParam(r, "ip")
	if net.ParseIP(ip) == nil {
		h.logger.Error().Err(fmt.Errorf("Invalid ip address: " + ip)).Msg("Invalid ip address")
		writeResponse(w, http.StatusBadRequest, model.Error{Error: "Internal server error"})
		return
	}
	host, err := h.service.ScanHost(ip)
	if err != nil {
		h.logger.Error().Err(err).Msg("ScanHost method error")
		writeResponse(w, http.StatusInternalServerError, model.Error{Error: "Internal server error"})
		return
	}
	writeResponse(w, http.StatusOK, host)
}
