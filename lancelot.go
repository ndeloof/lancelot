package lancelot

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	moby "github.com/docker/docker/client"
)

const version = "0.0.1"

func NewProxy(api moby.APIClient, cgroup string) Proxy {
	return &proxy{
		api:    api,
		cgroup: cgroup,
	}
}

type Proxy interface {
	Ping(w http.ResponseWriter, r *http.Request)
}

type proxy struct {
	api    moby.APIClient
	cgroup string
}

func (p *proxy) Ping(w http.ResponseWriter, r *http.Request) {
	ping, err := p.api.Ping(context.Background())
	if err != nil {
		panic("oups")
	}

	w.Header().Add("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Add("Pragma", "no-cache")
	bv := string(ping.BuilderVersion)
	if bv != "" {
		w.Header().Set("Builder-Version", bv)
	}
	w.Header().Set("Api-Version", ping.APIVersion)
	w.Header().Set("Docker-Experimental", strconv.FormatBool(ping.Experimental))
	w.Header().Set("Ostype", ping.OSType)
	w.Header().Set("Server", fmt.Sprintf("Lancelot/%s (%s)", version, ping.OSType))
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	if r.Method == http.MethodHead {
		w.Header().Set("Content-Length", "0")
		return
	}
	w.Write([]byte{'O', 'K'})
}
