package lancelot

import (
	"net"
	"net/http"
	"os"

	moby "github.com/docker/docker/client"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

const version = "0.0.1"

func NewProxy(api moby.APIClient, cgroup string) Proxy {
	return &proxy{
		api:    api,
		cgroup: cgroup,
	}
}

type Proxy struct {
	api    moby.APIClient
	cgroup string
}

func (p *Proxy) Serve(addr string) error {
	m := mux.NewRouter()
	p.RegisterRoutes(m)

	logs := handlers.LoggingHandler(os.Stdout, m)
	srv := &http.Server{Addr: addr, Handler: logs}
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	return srv.Serve(listener)
}

func (p *Proxy) RegisterRoutes(r *mux.Router) {
	r.Path("/_ping").Methods("GET").HandlerFunc(p.Ping)
}
