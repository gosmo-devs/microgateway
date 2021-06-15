package proxy

import (
	"net/http"
	"net/url"

	"github.com/gotway/gotway/internal/core"
	"github.com/gotway/gotway/pkg/log"
)

// Proxy interface
type Proxy interface {
	getTargetURL(r *http.Request) (*url.URL, error)
	ReverseProxy(w http.ResponseWriter, r *http.Request) error
}

// ResponseHandler is a function hook for handling responses
type ResponseHandler = func(serviceKey string, res *http.Response) error

type proxy struct {
	service        core.Service
	handleResponse ResponseHandler
	logger         log.Logger
}

func (p *proxy) log(req *http.Request, res *http.Response, target *url.URL) {
	p.logger.Infof("%s %s => %s %d", req.Method, req.URL, target, res.StatusCode)
}

func (p *proxy) handleError(w http.ResponseWriter, err error) {
	p.logger.Error("proxy error ", err)
	http.Error(w, "Internal server error", http.StatusInternalServerError)
}

func getDirector(target *url.URL) func(r *http.Request) {
	return func(r *http.Request) {
		r.Header.Add("X-Forwarded-Host", r.Host)
		r.Header.Add("X-Origin-Host", target.Host)
		r.URL.Scheme = target.Scheme
		r.URL.Host = target.Host
		r.URL.Path = target.Path
	}
}

// New instanciates a new Proxy
func New(service core.Service, handleResponse ResponseHandler, logger log.Logger) (Proxy, error) {
	proxy := proxy{
		service,
		handleResponse,
		logger,
	}
	switch service.Type {
	case core.ServiceTypeREST:
		return proxyREST{proxy}, nil
	case core.ServiceTypeGRPC:
		return proxyGRPC{proxy}, nil
	default:
		return nil, core.ErrInvalidServiceType
	}
}