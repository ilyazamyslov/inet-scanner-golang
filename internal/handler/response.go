package handler

import (
	"encoding/json"
	"net/http"

	"github.com/rs/zerolog"
)

func writeResponse(w http.ResponseWriter, code int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	b, err := json.Marshal(v)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"Internal server error"}`))
		return
	}
	w.WriteHeader(code)
	w.Write([]byte(b))
}

type Handler struct {
	logger  *zerolog.Logger
	service ScannerService
}

func New(logger *zerolog.Logger, srv ScannerService) *Handler {
	return &Handler{
		logger:  logger,
		service: srv,
	}
}
