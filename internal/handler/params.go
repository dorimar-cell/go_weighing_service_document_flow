package handler

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func parseNestedID(r *http.Request, name string) (int64, error) {
	return strconv.ParseInt(chi.URLParam(r, name), 10, 64)
}

func optionalQueryID(r *http.Request, key string) *int64 {
	value, err := strconv.ParseInt(r.URL.Query().Get(key), 10, 64)
	if err != nil || value <= 0 {
		return nil
	}
	return &value
}
