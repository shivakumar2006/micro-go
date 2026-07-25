package proxy

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"
)

type ServiceProxy struct {
	AuthProxy    *httputil.ReverseProxy
	CartProxy    *httputil.ReverseProxy
	VehicleProxy *httputil.ReverseProxy
}

func NewServiceProxy(authURL, cartURL, vehicleURL string) (*ServiceProxy, error) {
	authProxy, err := newReverseProxy(authURL)
	if err != nil {
		return nil, err
	}

	cartProxy, err := newReverseProxy(cartURL)
	if err != nil {
		return nil, err
	}

	vehicleProxy, err := newReverseProxy(vehicleURL)
	if err != nil {
		return nil, err
	}

	return &ServiceProxy{
		AuthProxy:    authProxy,
		CartProxy:    cartProxy,
		VehicleProxy: vehicleProxy,
	}, nil
}

func (s *ServiceProxy) Auth(w http.ResponseWriter, r *http.Request) {
	s.AuthProxy.ServeHTTP(w, r)
}

func (s *ServiceProxy) Cart(w http.ResponseWriter, r *http.Request) {
	s.CartProxy.ServeHTTP(w, r)
}

func (s *ServiceProxy) Vehicle(w http.ResponseWriter, r *http.Request) {
	s.VehicleProxy.ServeHTTP(w, r)
}

func newReverseProxy(targetURL string) (*httputil.ReverseProxy, error) {
	target, err := url.Parse(targetURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse URL: %w", err)
	}

	proxy := httputil.NewSingleHostReverseProxy(target)
	proxy.Transport = &http.Transport{
		MaxIdleConns:          100,
		MaxIdleConnsPerHost:   10,
		IdleConnTimeout:       90 * time.Second,
		ResponseHeaderTimeout: 30 * time.Second,
	}

	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		originalDirector(req)

		// clean headers we don't want to forward
		req.Header.Del("X-Forwarded-For")
		req.Header.Del("CP-Connected-IP")

		// set forwarded headers
		req.Header.Set("X-Forwarded-Host", req.Host)
		req.Header.Set("X-Origin-Host", target.Host)

		req.Header.Set("X-Gateway", "vehicle-management-system-gateway")
	}

	proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadGateway)
		w.Write([]byte(`{
			"error": "downstream service unreachable",
			"detail": "` + err.Error() + `"
		}`))
	}

	return proxy, nil
}
