package http

import (
	"encoding/json"
	"net/http"

	"github.com/Alma-media/kuzya/state"
)

type StateReader interface {
	List() ([]state.Device, error)
}

type Handler struct {
	state StateReader
}

func NewHandler(state StateReader) *Handler {
	return &Handler{state: state}
}

func (h *Handler) StatusHandler(w http.ResponseWriter, r *http.Request) {
	devices, err := h.state.List()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	if devices == nil {
		devices = []state.Device{}
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(devices); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}
}
