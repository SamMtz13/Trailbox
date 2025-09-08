package http

import (
	"fmt"
	"log"
	"net/http"
)

func Healthz(w http.ResponseWriter, r *http.Request) {
	log.Printf("[routes] health called from %s", r.RemoteAddr)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprintf(w, "%s service ok", "routes")
}
