package helpers

import (
	"net/url"
)

//Validates the URL is OK
func IsValidUrl(urlString string) bool {
	if urlString == "" {
		return false
	}
	_, err := url.Parse(urlString)
	if err != nil {
		return false
	}
	return true
}
