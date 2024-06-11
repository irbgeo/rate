package http

import (
	"net/http"

	"github.com/gorilla/mux"
)

var (
	rateURL    = "/rate"
	rateMethod = http.MethodGet
)

type rateServer struct {
	svc svc

	transport *rateServerTransport
}

func rateHandle(router *mux.Router, svc svc) {
	s := rateServer{
		svc:       svc,
		transport: &rateServerTransport{},
	}
	router.HandleFunc(rateURL, s.Rate).Methods(rateMethod)
}

func (s *rateServer) Rate(w http.ResponseWriter, r *http.Request) {
	req, err := s.transport.decodeRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp, err := s.svc.Rate(r.Context(), *req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	s.transport.encodeResponse(w, resp)
}
