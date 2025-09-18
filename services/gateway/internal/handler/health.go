package handler

import (
	"fmt"
	"net/http"
)

// Healthz handler
func Healthz(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "OK - gateway")
}
