package handler

import "net/http"

// HealthHandler provides a health check endpoint
func HealthHandler(w http.ResponseWriter, r *http.Request) {
	msg := `{"status": "ok"}`

	w.WriteHeader(http.StatusOK)
	//nolint:errcheck
	w.Write([]byte(msg))
}
