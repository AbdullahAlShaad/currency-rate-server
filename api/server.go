package api

import (
	"github.com/shaad7/currency-rate-server/controller"
	"github.com/shaad7/currency-rate-server/util"
	"log"
	"net/http"
	"strings"
)

// Server is the HTTP api.
type Server struct {
	listenAddress string
	controller    *controller.Controller
}

// NewServer creates a new api using the listenAddress and controller information
func NewServer(listenAddress string, controller *controller.Controller) *Server {
	return &Server{
		listenAddress: listenAddress,
		controller:    controller,
	}
}

// Start function starts the API Server. It registers two endpoint and call a generalize
// function for all others endpoint
func (s *Server) Start() error {
	http.HandleFunc("/rates/analyze", s.analyzeRate)
	http.HandleFunc("/rates/latest", s.getLatestRate)
	http.HandleFunc("/", s.ServeURLs)

	log.Printf("Listenning to the address : %s", s.listenAddress)
	return http.ListenAndServe(s.listenAddress, nil)
}

// ServeURLs function handles all the endpoint except /rates/analyze and /rates/latest. First the verbs except get is rejected.
// Then we match if the endpoint has /rate/analyze/yyyy-mm-dd format, if yes it is handled. Otherwise, error thrown/
func (s *Server) ServeURLs(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		util.JsonError(w, http.StatusMethodNotAllowed, util.ErrorMethodNotAllowed, nil)
	}
	if util.DateAPIPattern.MatchString(r.URL.Path) {
		date := strings.TrimPrefix(r.URL.Path, "/rates/")
		log.Printf("Get Rates of Date : %s", date)
		s.getRate(date, w, r)
	} else {
		log.Printf("Invalid API path %s\n", r.URL)
		util.JsonError(w, http.StatusInternalServerError, util.ErrorForbidden, nil)
	}
}

// analyzeRate functions handles /rates/analyze endpoint for GET verb.
func (s *Server) analyzeRate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		util.JsonError(w, http.StatusMethodNotAllowed, util.ErrorMethodNotAllowed, nil)
	}
	log.Println("Handling /rates/analyze endpoint")
	analysis, err := s.controller.GetAnalysis()
	if err != nil {
		data := map[string]interface{}{
			"message": err.Error(),
		}
		util.JsonError(w, http.StatusInternalServerError, util.ErrorNotFound, data)
	}
	util.WriteJSONResponse(w, http.StatusOK, analysis)
}

// getLatestRate handles /rates/latest endpoint for GET verb
func (s *Server) getLatestRate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		util.JsonError(w, http.StatusMethodNotAllowed, util.ErrorMethodNotAllowed, nil)
	}
	log.Println("Handling /rates/latest endpoint")
	latestRate, err := s.controller.GetLatestExchangeRate()
	if err != nil {
		data := map[string]interface{}{
			"message": err.Error(),
		}
		util.JsonError(w, http.StatusInternalServerError, util.ErrorNotFound, data)
	}
	util.WriteJSONResponse(w, http.StatusOK, latestRate)
}

// getRate handles /rates/yyyy-mm-dd format endpoints for GET verb
func (s *Server) getRate(date string, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		util.JsonError(w, http.StatusMethodNotAllowed, util.ErrorMethodNotAllowed, nil)
	}
	exchangeRate, err := s.controller.GetExchangeRate(date)
	if err != nil {
		data := map[string]interface{}{
			"message": err.Error(),
		}
		util.JsonError(w, http.StatusInternalServerError, util.ErrorNotFound, data)
	}
	util.WriteJSONResponse(w, http.StatusOK, exchangeRate)
}
