package handler

import (
	"github.com/gorilla/mux"
	"github.com/jaegertracing/jaeger/model"
	"net/http"
)

type ApiHandler struct {
}

func NewApiHandler() *ApiHandler {
	return &ApiHandler{}
}

const (
	TraceIDParam = "traceID"
)

func (h *ApiHandler) RegisterRouter(r *mux.Router) {
	r.HandleFunc("/traces/{traceID}", h.getTrace).Methods(http.MethodGet)
	r.HandleFunc("/archive/{traceID}", h.archiveTrace).Methods(http.MethodGet)
	r.HandleFunc("/traces", h.search).Methods(http.MethodGet)
	r.HandleFunc("/services", h.getServices).Methods(http.MethodGet)
	r.HandleFunc("/operations", h.getOperations).Methods(http.MethodGet)
	r.HandleFunc("/service/{service}/operations", h.getOperationsLegacy).Methods(http.MethodGet)
	r.HandleFunc("/dependencies", h.dependencies).Methods(http.MethodGet)
}

func (h *ApiHandler) parseTraceId(writer http.ResponseWriter, r *http.Request) {
	traceIDVar := mux.Vars(r)["traceID"]
	traceID, err := model.TraceIDFromString(traceIDVar)
	if h.handlerError(writer, err, http.StatusNotFound) {
		//return traceid
	}
}
func (h *ApiHandler) getTrace(writer http.ResponseWriter, request *http.Request) {
	traceID := mux.Vars(request)["traceID"]
	h.parseTraceId(writer, request)
}

func (h *ApiHandler) handlerError(writer http.ResponseWriter, err error, statusCode int) bool {
	if err == nil {
		return false
	}

	return true
}

func (h *ApiHandler) archiveTrace(writer http.ResponseWriter, request *http.Request) {

}

func (h *ApiHandler) search(writer http.ResponseWriter, request *http.Request) {

}

func (h *ApiHandler) getServices(writer http.ResponseWriter, request *http.Request) {

}

func (h *ApiHandler) getOperations(writer http.ResponseWriter, request *http.Request) {

}

func (h *ApiHandler) getOperationsLegacy(writer http.ResponseWriter, request *http.Request) {

}

func (h *ApiHandler) dependencies(writer http.ResponseWriter, request *http.Request) {

}
