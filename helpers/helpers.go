package helpers

import (
	"encoding/json"
	"fmt"
	"net/http"
)


func RespondWithError(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	err := json.NewEncoder(w).Encode(map[string]string{
		"error": message,
	})
	
	if err != nil {
		http.Error(w, fmt.Sprintf("%s", err), code)
		return;
	}
}