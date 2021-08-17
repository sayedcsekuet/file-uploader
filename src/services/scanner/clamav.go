package scanner

import (
	"github.com/dutchcoders/go-clamd"
	log "github.com/sirupsen/logrus"
	"io"
)

//Clamav scans files using clamav
type Clamav struct {
	Scanner
	address string
	debug   bool
	clam    *clamd.Clamd
}

//Returns the engine address.
func (c *Clamav) Address() string {
	return c.address
}
func (c *Clamav) SetAddress(url string) {
	c.clam = clamd.NewClamd(url)
	if c.debug {
		log.Info("Initialised clamav connection to "+url, nil)
	}
	c.address = url
}

func (c *Clamav) HasVirus(reader io.Reader) (bool, error) {
	result, err := c.Scan(reader)
	if err != nil {
		return false, err
	}

	return result.Virus, nil
}

func (c *Clamav) Scan(reader io.Reader) (*Result, error) {
	if c.debug {
		log.Info("Sending to clamav", nil)
	}

	ch, err := c.clam.ScanStream(reader, nil)
	if err != nil {
		return nil, err
	}
	var status string

	r := <-ch

	switch r.Status {
	case clamd.RES_OK:
		status = RES_CLEAN
	case clamd.RES_FOUND:
		status = RES_FOUND
	case clamd.RES_ERROR:
	case clamd.RES_PARSE_ERROR:
	default:
		status = RES_ERROR
	}

	result := &Result{
		Status:      status,
		Virus:       status == RES_FOUND,
		Description: r.Description,
	}

	if c.debug {
		log.Info(result, nil)
	}

	return result, nil
}

func (c *Clamav) Ping() error {
	return c.clam.Ping()
}

func (c *Clamav) Version() (string, error) {
	ch, err := c.clam.Version()
	if err != nil {
		return "", err
	}

	r := <-ch
	return r.Raw, nil
}
