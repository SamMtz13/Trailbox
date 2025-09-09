package http

import (
	"context"
	"net/http"
	"time"

	"trailbox/services/routes/internal/client/peer"
)

// PeerHandler expone endpoints que llaman a otro microservicio ("curl" en Go).
type PeerHandler struct {
	client *peer.Client
}

// NewPeerHandler construye el handler. Si PEER_URL no está seteada, devolverá 503.
func NewPeerHandler() *PeerHandler {
	c, err := peer.New()
	if err != nil {
		return &PeerHandler{client: nil}
	}
	return &PeerHandler{client: c}
}

// GET /peer/health -> llama GET {PEER_URL}/health y devuelve la respuesta
func (h *PeerHandler) HandlePeerHealth(w http.ResponseWriter, r *http.Request) {
	if h.client == nil {
		http.Error(w, "peer not configured (set PEER_URL)", http.StatusServiceUnavailable)
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	code, body, err := h.client.Get(ctx, "/health")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(code)
	_, _ = w.Write(body)
}

// GET /peer/proxy?path=/v1/lo-que-sea -> reenvía GET al peer en ese path
func (h *PeerHandler) HandlePeerProxy(w http.ResponseWriter, r *http.Request) {
	if h.client == nil {
		http.Error(w, "peer not configured (set PEER_URL)", http.StatusServiceUnavailable)
		return
	}
	path := r.URL.Query().Get("path")
	if path == "" {
		http.Error(w, "missing query param: path", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 8*time.Second)
	defer cancel()

	code, body, err := h.client.Get(ctx, path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
	w.Header().Set("Content-Type", "application/octet-stream")
	w.WriteHeader(code)
	_, _ = w.Write(body)
}

// Heartbeat pega periódicamente al peer para revisar salud.
type Heartbeat struct {
	h      *PeerHandler
	target string
	period time.Duration
}

// NewHeartbeat crea un heartbeat que consulta GET {PEER_URL}{target} cada periodo.
func NewHeartbeat(h *PeerHandler, target string, period time.Duration) *Heartbeat {
	return &Heartbeat{h: h, target: target, period: period}
}

func (hb *Heartbeat) Run(stop <-chan struct{}, logf func(string, ...any)) {
	if hb.h == nil || hb.h.client == nil {
		logf("[routes] heartbeat disabled: peer not configured")
		return
	}
	t := time.NewTicker(hb.period)
	defer t.Stop()
	logf("[routes] heartbeat started: target=%s interval=%s", hb.target, hb.period)

	for {
		select {
		case <-t.C:
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			_, _, err := hb.h.client.Get(ctx, hb.target)
			cancel()
			if err != nil {
				logf("[routes] heartbeat FAIL: %v", err)
				continue
			}
			logf("[routes] heartbeat OK")
		case <-stop:
			logf("[routes] heartbeat stopped")
			return
		}
	}
}
