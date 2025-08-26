// Package http is the input adapter for the users entries
package http

import (
	"log"
	"net/http"

	"github.com/FreyreCorona/Shortly/src/redirect_svc/internal/application"
)

type Handler struct {
	service *application.RedirectionService
}

func NewHandler(urlService *application.RedirectionService) *Handler {
	return &Handler{service: urlService}
}

func (h *Handler) Routes(mux *http.ServeMux) {
	mux.HandleFunc("/shortly/", h.redirect)
}

func (h *Handler) redirect(w http.ResponseWriter, r *http.Request) {
	short := r.URL.Path[1:] // gets all from http://->
	if short == "" {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	url, err := h.service.GetURL(short)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Printf("error retrieving url %s", err.Error())
		return
	}

	http.Redirect(w, r, url.RawURL, http.StatusFound)
}
