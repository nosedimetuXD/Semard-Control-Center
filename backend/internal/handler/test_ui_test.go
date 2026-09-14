package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestServeTestUI(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "/test", nil)
	if err != nil {
		t.Fatalf("Error al crear request: %v", err)
	}

	rr := httptest.NewRecorder()
	ServeTestUI(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Se esperaba código %d, se obtuvo %d", http.StatusOK, rr.Code)
	}

	contentType := rr.Header().Get("Content-Type")
	if !strings.Contains(contentType, "text/html") {
		t.Errorf("Se esperaba Content-Type text/html, se obtuvo %s", contentType)
	}

	body := rr.Body.String()
	if !strings.Contains(body, "SEMARD Control Center") {
		t.Errorf("El cuerpo de la respuesta no contiene el título esperado")
	}
	if !strings.Contains(body, "Test Studio") {
		t.Errorf("El cuerpo no contiene 'Test Studio'")
	}
}
