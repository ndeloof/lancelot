package lancelot

import (
	"context"
	"net/http"

	"github.com/docker/docker/api/server/httputils"
	moby "github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/gorilla/mux"
)

func (p *Proxy) containerList(w http.ResponseWriter, r *http.Request) {
	if err := httputils.ParseForm(r); err != nil {
		p.error(w, err)
		return
	}
	filter, err := filters.FromJSON(r.Form.Get("filters"))
	if err != nil {
		p.error(w, err)
		return
	}
	filter.Add("label", "com.docker.lancelot="+p.cgroup)

	config := moby.ContainerListOptions{
		All:     httputils.BoolValue(r, "all"),
		Size:    httputils.BoolValue(r, "size"),
		Since:   r.Form.Get("since"),
		Before:  r.Form.Get("before"),
		Filters: filter,
	}

	containers, err := p.api.ContainerList(context.Background(), config)
	if err != nil {
		p.error(w, err)
		return
	}
	for i, c := range containers {
		delete(c.Labels, "com.docker.lancelot")
		containers[i] = c
	}

	httputils.WriteJSON(w, http.StatusOK, containers)
}

func (p *Proxy) containerInspect(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]
	container, err := p.api.ContainerInspect(context.Background(), name)
	if err != nil {
		p.error(w, err)
		return
	}
	if _, ok := container.Config.Labels["com.docker.lancelot"]; !ok {
		p.error(w, objNotFoundError{"container", name})
		return
	}

	delete(container.Config.Labels, "com.docker.lancelot")
	httputils.WriteJSON(w, http.StatusOK, container)
}
