package lancelot

import (
	"context"
	"net/http"

	"github.com/docker/docker/api/server/httputils"
	moby "github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
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

	httputils.WriteJSON(w, http.StatusOK, containers)
}
