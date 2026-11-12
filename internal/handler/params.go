package handler

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func parseNestedID(r *http.Request, name string) (int64, error) {
	return strconv.ParseInt(chi.URLParam(r, name), 10, 64)
}
