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
	var netErr net.Error
	if errors.As(err, &netErr) && netErr.Timeout() {
		return true
	}

	if errors.As(err, &netErr) {
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
