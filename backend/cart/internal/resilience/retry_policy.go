package resilience

import (
	"errors"
	"net"
	"net/http"
)

type HTTPStatusError struct {
	StatusCode int
}

func (e *HTTPStatusError) Error() string {
	return http.StatusText(e.StatusCode)
}

func IsRetryable(err error) bool {

	// network timeout - network error
	var netError net.Error
	if errors.As(err, &netError) && netError.Timeout() {
		return true
	}

	// connection refused or temporay network error
	if errors.As(err, &netError) {
		return true
	}

	// http status error
	var statusErr *HTTPStatusError

	if errors.As(err, &statusErr) {
		switch statusErr.StatusCode {
		case http.StatusBadGateway, http.StatusServiceUnavailable, http.StatusGatewayTimeout:
			return true
		}
	}

	return false
}
