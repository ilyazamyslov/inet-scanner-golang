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
	service Service
}

type Service interface {
	//ScanByIp(string, time.Time) (*model.Price, error)
}

func New(logger *zerolog.Logger, srv Service) *Handler {
	return &Handler{
		logger:  logger,
		service: srv,
	}
}
