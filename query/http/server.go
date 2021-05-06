package http

import (
	"github.com/gorilla/mux"
	"net/http"
	"query/http/handler"
)

func createHTTPServer() *http.Server {

	apiHandler := handler.NewApiHandler()

	r := mux.NewRouter()

	apiHandler.RegisterRouter(r)

	http.Handle("/", r)
	return &http.Server{
		Handler: r,
	}
}
