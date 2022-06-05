package helpers

import (
	"errors"
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

