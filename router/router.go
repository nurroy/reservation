package router

import (
	pg "belajar/reservation/pkg/v1/postgres"
	"github.com/kenshaw/envcfg"
	"github.com/sirupsen/logrus"
	"github.com/soheilhy/cmux"
	"log"
	"net"
	"net/http"
)

var gconfig *envcfg.Envcfg

// Setup sets up the environment.
func Setup() (*envcfg.Envcfg, *logrus.Logger, error) {
	// create config and logger
	config, err := envcfg.New()
	if err != nil {
		return nil, nil, err
	}
	logger := logrus.New()

	// force all writes to regular log to logger
	log.SetOutput(logger.Writer())
	log.SetFlags(0)

	// configure logging for environment
	logger.Formatter = &logrus.TextFormatter{
		ForceColors:   true,
		ForceQuote:    true,
		FullTimestamp: true,
	}

	logger.Info("[CONFIG] Setup complete")


	// setup postgres
	err = pg.New(config)
	if err != nil {
		return nil, nil, err
	}

	logger.Infof("[POSTGRES] Setup complete")


	// assign config become global
	gconfig = config

	return config, logger, nil
}

// IgnoreErr returns true when err can be safely ignored.
func IgnoreErr(err error) bool {
	switch {
	case err == nil || err == http.ErrServerClosed || err == cmux.ErrListenerClosed:
		return true
	}
	if opErr, ok := err.(*net.OpError); ok {
		return opErr.Err.Error() == "use of closed network connection"
	}
	return false
}

