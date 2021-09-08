package handler

import (
	"fmt"
	"net"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/ilyazamyslov/inet-scanner-golang/internal/model"
	"github.com/rs/zerolog"
)

const (
	ScanNetworkPath = "/network/{ip}/{lenMask}"
)

type NetwokScan struct {
	logger  *zerolog.Logger
	service ScannerService
}

type ScannerService interface {
	ScanNetwork(ip string, lenMask int) ([]model.Host, error)
	ScanHost(ip string) (model.Host, error)
}

func NewNetworkScan(logger *zerolog.Logger, service ScannerService) *NetwokScan {
	return &NetwokScan{
		logger:  logger,
		service: service,
	}
}

func (h *NetwokScan) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ip := chi.URLParam(r, "ip")
	if net.ParseIP(ip) == nil {
		h.logger.Error().Err(fmt.Errorf("Invalid ip address: " + ip)).Msg("Invalid ip address")
		writeResponse(w, http.StatusBadRequest, model.Error{Error: "Bad request"})
		return
	}
	lenMask, err := strconv.Atoi(chi.URLParam(r, "lenMask"))
	if err != nil {
		h.logger.Error().Err(err).Msg("Invalid mask")
		writeResponse(w, http.StatusBadRequest, model.Error{Error: "Bad request"})
		return
	}
	if lenMask < 0 || lenMask > 32 {
		h.logger.Error().Err(fmt.Errorf("Invalid mask: " + ip)).Msg("Invalid mask")
		writeResponse(w, http.StatusBadRequest, model.Error{Error: "Bad request"})
		return
	}
	hosts, err := h.service.ScanNetwork(ip, lenMask)
	if err != nil {
		h.logger.Error().Err(err).Msg("ScanHost method error")
		writeResponse(w, http.StatusInternalServerError, model.Error{Error: "Internal server error"})
		return
	}
	writeResponse(w, http.StatusOK, hosts)
}
