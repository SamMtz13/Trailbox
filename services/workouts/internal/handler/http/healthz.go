package http

import (
	"fmt"
	"log"
	"net/http"
)

func Healthz(w http.ResponseWriter, r *http.Request) {
	log.Printf("[workouts] healthz called from %s", r.RemoteAddr)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprintf(w, "%s service ok", "workouts")
}
