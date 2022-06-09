package helpers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	schema "github.com/gorilla/schema"
)

func ParseForm(r *http.Request, dst interface{}) error {
	err := r.ParseForm()
	if err != nil {
		return errors.New("error while parsing form")
	}
	decoder := schema.NewDecoder()
	err = decoder.Decode(dst, r.PostForm)
	if err != nil {
		return errors.New("bad data received")
	}
	return nil
}

func RespondWithError(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	err := json.NewEncoder(w).Encode(map[string]string{
		"error": message,
		"code": fmt.Sprintf("%d", code),
	})
	
	if err != nil {
		http.Error(w, fmt.Sprintf("%s", err), code)
		return
	}
}