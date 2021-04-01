package menu

import (
	"github.com/kenshaw/envcfg"
	"github.com/sirupsen/logrus"
)

// Methode is the methode type
type Method int

const (
	// List of different Methods
	Get Method = iota
	POST
	other
)


// Server is the server object for this api service.
type Server struct {
	config *envcfg.Envcfg
	logger *logrus.Logger
}

//RetrieveRequest
type RetrieveRequest struct {
	Id string `json:"id"`
	Name string `json:"name"`
	Description string `'json:"desc"'`
}

//RetrieveRequest
type RetrieveResponse struct {
	Id string `json:"id"`
	Name string `json:"name"`
	Description string `'json:"desc"'`
}


// New creates a new server.
func New(config *envcfg.Envcfg, logger *logrus.Logger) *Server {
	return &Server{
		config: config,
		logger: logger,
	}
}

