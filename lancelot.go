package lancelot

import (
	"fmt"
	"net"
	"net/http"
	"os"

	"github.com/docker/docker/api/server/httputils"
	"github.com/docker/docker/api/types"
	moby "github.com/docker/docker/client"
	"github.com/docker/docker/errdefs"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

const version = "0.0.1"

func NewProxy(api moby.APIClient, cgroup string) *Proxy {
	return &Proxy{
		api:    api,
		cgroup: cgroup,
	}
}

type Proxy struct {
	api        moby.APIClient
	cgroup     string
	containers []string
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
	r.Path("/_ping").Methods("GET").HandlerFunc(p.ping)
	r.Path("/_ping").Methods("HEAD").HandlerFunc(p.ping)
	r.Path("/version").Methods("GET").HandlerFunc(p.version)
	r.Path("/v{version:[0-9.]+}/version").Methods("GET").HandlerFunc(p.version)
	r.Path("/info").Methods("GET").HandlerFunc(p.info)
	r.Path("/v{version:[0-9.]+}/info").Methods("GET").HandlerFunc(p.info)

	r.Path("/v{version:[0-9.]+}/containers/json").Methods("GET").HandlerFunc(p.containerList)
	r.Path("/v{version:[0-9.]+}/containers/{name:.*}/json").Methods("GET").HandlerFunc(p.containerInspect)
	r.Path("/v{version:[0-9.]+}/containers/create").Methods("POST").HandlerFunc(p.containerCreate)
}

func (p *Proxy) error(w http.ResponseWriter, err error) {
	fmt.Fprintln(os.Stderr, err.Error())
	response := &types.ErrorResponse{
		Message: err.Error(),
	}
	httputils.WriteJSON(w, errdefs.GetHTTPErrorStatusCode(err), response)
}

type objNotFoundError struct {
	object string
	id     string
}

func (e objNotFoundError) Error() string {
	return "No such " + e.object + ": " + e.id
}

func (e objNotFoundError) NotFound() {}
