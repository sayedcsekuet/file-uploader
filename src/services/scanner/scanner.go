package scanner

import (
	"fmt"
	"io"
)

//Scanner result statuses
const (
	RES_CLEAN = "CLEAN"
	RES_FOUND = "FOUND"
	RES_ERROR = "ERROR"
)

/*
 The Scanner interface. This is meant as a starting point to support multiple
 virus scanners. The concrete struct type is the Engine, that should then be
 embedded in the scanner implementation. See scanner/clamav.go for the current
 and only Scanner implementation.
*/
type Scanner interface {
	//Sets the scanner engine address
	SetAddress(address string)
	//Gets the scanner engine address
	Address() string
	//This function performs the actual virus scan and returns a boolean indicating whether a
	//virus has been found or not
	HasVirus(reader io.Reader) (bool, error)
	//This function performs the actual virus scan and returns an engine-specific response string
	Scan(reader io.Reader) (*Result, error)
	//Tests the liveliness of the underlying scan engine
	Ping() error
	//Returns the version of the underlying scan engine
	Version() (string, error)
}

/*
Embeds a scan result.Status is one of the RES constants
Virus is true or false depending a Virus has been detected
Description is an extended status, containing the virus name
*/
type Result struct {
	Status      string
	Virus       bool
	Description string
}

func (r *Result) String() string {
	ret := fmt.Sprintf("Status: %s; Virus: %v", r.Status, r.Virus)

	if r.Virus {
		ret += fmt.Sprintf("; Description: %s", r.Description)
	}

	return ret
}
