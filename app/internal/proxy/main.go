package proxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

func NewProxy(targetHost string) (*httputil.ReverseProxy, error) {
	u, err := url.Parse(targetHost)
	if err != nil {
		return nil, err
	}

	return httputil.NewSingleHostReverseProxy(u), nil
}

func ProxyRequestHandler(proxyPrefix string, proxy *httputil.ReverseProxy) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		r.URL.Path = strings.TrimLeft(path, proxyPrefix)
		proxy.ServeHTTP(w, r)
	}
}
