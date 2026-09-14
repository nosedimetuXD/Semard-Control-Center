package handler

import (
	_ "embed"
	"net/http"
)

//go:embed test.html
var testHTML []byte

// ServeTestUI sirve la interfaz web interactiva de pruebas de la API embebida directamente en el binario Go.
func ServeTestUI(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(testHTML)
}
