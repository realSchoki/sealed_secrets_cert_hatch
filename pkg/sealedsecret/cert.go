package sealedsecret

import (
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"os"
	"time"

	"github.com/go-co-op/gocron/v2"
	log "github.com/sirupsen/logrus"
)

type Cert struct {
	path    string
	content []byte
	hash    [32]byte
}

func (c *Cert) scheduleRegularRefresh() (*gocron.Scheduler, error) {
	s, err := gocron.NewScheduler()
	if err != nil {
		return nil, err
	}

	_, err = s.NewJob(
		gocron.DurationJob(10*time.Second),
		gocron.NewTask(c.load),
	)
	if err != nil {
		return nil, err
	}

	s.Start()
	return &s, nil
}

func (c *Cert) load() error {
	log.Debug("Reading certfile")
	bytes, err := os.ReadFile(c.path)
	if err != nil {
		log.Debug("Error while reading certfile")
		if c.Content() != nil {
			log.Error("Error reading cert, using old cert")
		}
		return err
	} else {
		log.Debugf("Cert read: %v", c.path)
	}

	sha256sum := sha256.Sum256(bytes)
	if c.hash == sha256sum {
		log.Debug("Cert did not change")
		return nil
	} else {
		log.Info("Cert changed")
		c.hash = sha256sum
	}

	block, _ := pem.Decode(bytes)
	if err != nil {
		if c.Content() != nil {
			log.Error("Error reading cert, using old cert")
		}
		log.Errorf("Error decoding cert: %v", err)
		return err
	}

	parsedCert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		if c.Content() != nil {
			log.Error("Error reading cert, using old cert")
		}
		log.Errorf("Error parsing cert: %v", err)
		return err
	} else {
		log.Infof("Cert loaded: %v", c.path)
		log.Infof("Valid until: %v", parsedCert.NotAfter)
	}

	c.content = bytes
	return nil
}

func (c *Cert) Content() []byte {
	return c.content
}

func NewCert(path string) (*Cert, *gocron.Scheduler, error) {
	cert := Cert{
		path: path,
	}

	err := cert.load()
	if err != nil {
		return nil, nil, err
	}
	scheduler, err := cert.scheduleRegularRefresh()
	if err != nil {
		log.Error("Error scheduling cert refresh")
		return nil, nil, err
	}

	return &cert, scheduler, nil
}
