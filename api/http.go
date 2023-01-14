package api

import (
	"io"
	"log"
	"net/http"

	"github.com/pkg/errors"

	"github.com/dimoynwa/url-shortener/serializer/json"
	"github.com/dimoynwa/url-shortener/serializer/msgpack"
	"github.com/dimoynwa/url-shortener/shortener"
	"github.com/go-chi/chi/v5"
)

type RedirectHandler interface {
	Get(http.ResponseWriter, *http.Request)
	Post(http.ResponseWriter, *http.Request)
}

type handler struct {
	redirectService shortener.RedirectService
}

func NewHandler(service shortener.RedirectService) RedirectHandler {
	return &handler{redirectService: service}
}

func setUpResponse(wr http.ResponseWriter, contentType string, statusCode int, body []byte) {
	wr.Header().Set("Content-Type", contentType)
	wr.WriteHeader(statusCode)

	_, err := wr.Write(body)
	if err != nil {
		log.Println(err)
	}
}

func (h *handler) serializer(contentType string) shortener.RedirectSerializer {
	if contentType == "application/x-msgpack" {
		return &msgpack.Redirect{}
	}
	return &json.Redirect{}
}

func (h *handler) Get(writer http.ResponseWriter, request *http.Request) {
	code := chi.URLParam(request, "code")

	redirect, err := h.redirectService.Find(code)
	if err != nil {
		if errors.Cause(err) == shortener.ErrRedirectNotFount {
			http.Error(writer, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	http.Redirect(writer, request, redirect.URL, http.StatusMovedPermanently)

}

func (h *handler) Post(writer http.ResponseWriter, request *http.Request) {
	contentType := request.Header.Get("Content-Type")
	body, err := io.ReadAll(request.Body)
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	serializer := h.serializer(contentType)

	redirect, err := serializer.Decode(body)
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	err = h.redirectService.Store(redirect)
	if err != nil {
		if errors.Cause(err) == shortener.ErrRedirectInvalid {
			http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	responseBody, err := serializer.Encode(redirect)
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	setUpResponse(writer, contentType, http.StatusCreated, responseBody)
}
