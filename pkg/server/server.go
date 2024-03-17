package server

import (
	"errors"
	"net/http"
	"os"

	"github.com/realSchoki/sealed_secrets_cert_hatch/pkg/sealedsecret"

	log "github.com/sirupsen/logrus"
)

func StartServer() {
	log.SetFormatter(&log.JSONFormatter{})
	if os.Getenv("DEBUG") == "true" {
		log.SetLevel(log.DebugLevel)
		log.Debug("Debug logging activated")
	}

	path, exists := os.LookupEnv("CERT_PATH")
	if !exists {
		log.Panic("CERT_PATH not set")
	} else {
		log.Infof("Certfile path: %s", path)
	}

	cert, scheduler, err := sealedsecret.NewCert(path)
	if err != nil {
		log.Panic("Cannot load certfile")
	}

	server := http.NewServeMux()
	server.Handle("/cert", http.HandlerFunc(getCertHandler(cert)))

	err = http.ListenAndServe(":8080", server)
	if err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			log.Info("Server closed")
		} else {
			log.Infof("Error starting server: %v", err)
		}
	}

	log.Info("Server stopped")
	(*scheduler).Shutdown()
}

func getCertHandler(cert *sealedsecret.Cert) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if cert.Content() == nil {
			http.Error(w, "Error reading cert", http.StatusInternalServerError)
		} else {
			w.Write(cert.Content())
		}
	}
}
