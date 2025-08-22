// Package http is the input adapter for the users entries
package http

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/FreyreCorona/Shortly/url_shortener_svc/internal/application"
)

type Handler struct {
	urlService *application.CreateURLService
}

type JSONBody map[string]any

func NewHandler(urlService *application.CreateURLService) *Handler {
	return &Handler{urlService: urlService}
}

func (h *Handler) Routes(mux *http.ServeMux) {
	mux.HandleFunc("/shortly", h.createURL)
	mux.HandleFunc("/", h.redirect)
}

func (h *Handler) createURL(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	var input struct {
		RawURL string `json:"raw_url"`
	}
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	url, err := h.urlService.CreateURL(input.RawURL)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Printf("error on createURl: %s", err.Error())
		return
	}

	err = json.NewEncoder(w).Encode(JSONBody{"id": url.ID, "raw_url": url.RawURL, "short_code": url.ShortCode, "created_at": url.CreatedAt})
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Printf("error sending a response :%s", err.Error())
	}
}

func (h *Handler) redirect(w http.ResponseWriter, r *http.Request) {
	short := r.URL.Path[1:] // gets all from http://->
	if short == "" {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	url, err := h.urlService.GetURL(short)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Printf("error retrieving url %s", err.Error())
		return
	}

	http.Redirect(w, r, url.RawURL, http.StatusFound)
}
