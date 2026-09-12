package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHubInfo(t *testing.T) {
	hubHandler := NewHubHandler(nil)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/hub/info", nil)
	w := httptest.NewRecorder()

	hubHandler.GetHubInfo(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Esperaba status 200, obtuvo %d", w.Code)
	}

	var response map[string]interface{}
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("Error al decodificar JSON: %v", err)
	}

	if response["acronym"] != "SEMARD" {
		t.Errorf("Esperaba acronym SEMARD, obtuvo %v", response["acronym"])
	}

	if response["institution"] != "Universidad de Cartagena" {
		t.Errorf("Esperaba institución Universidad de Cartagena, obtuvo %v", response["institution"])
	}

	lines, ok := response["research_lines"].([]interface{})
	if !ok || len(lines) == 0 {
		t.Errorf("Esperaba líneas de investigación no vacías, obtuvo %v", response["research_lines"])
	}
}
