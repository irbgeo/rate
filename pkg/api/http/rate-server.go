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

//	@Summary		Get exchange rate
//	@Description	Returns the exchange rate from tokenIn to tokenOut, which are taken from the request parameters
//	@Tags			rate
//	@Accept			json
//	@Produce		json
//	@Param			tokenIn		query		string	true	"Source currency"
//	@Param			tokenOut	query		string	true	"Target currency"
//	@Success		200			{object}	string	"Successful response with the exchange rate"
//	@Router			/rate [get]
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
