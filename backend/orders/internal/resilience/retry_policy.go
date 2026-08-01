package resilience

import (
	"errors"
	"net"
	"net/http"
)

type HTTPStatusError struct {
	StatusCode int
}

func (h HTTPStatusError) Error() string {
	return http.StatusText(h.StatusCode)
}

func IsRetyrable(err error) bool {
	var netError net.Error

	if errors.As(err, &netError) && netError.Timeout() {
		return true
	}

	if errors.As(err, &netError); err != nil {
		return true
	}

	var statusErr *HTTPStatusError

	if errors.As(err, &statusErr) {
		switch statusErr.StatusCode {
		case http.StatusBadGateway, http.StatusServiceUnavailable, http.StatusGatewayTimeout:
			return true
		}
	}

	return false
}
