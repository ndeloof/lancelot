package lancelot

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/docker/docker/api/server/httputils"
	moby "github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	networktypes "github.com/docker/docker/api/types/network"
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

// FIXME get some dependency issue when I use runconfig.ContainerConfigWrapper
type ContainerConfigWrapper struct {
	*container.Config
	HostConfig       *container.HostConfig
	NetworkingConfig *networktypes.NetworkingConfig `json:"NetworkingConfig,omitempty"`
}

func (p *Proxy) containerCreate(w http.ResponseWriter, r *http.Request) {
	if err := httputils.ParseForm(r); err != nil {
		p.error(w, err)
		return
	}
	if err := httputils.CheckForJSON(r); err != nil {
		p.error(w, err)
		return
	}

	name := r.Form.Get("name")

	var cc ContainerConfigWrapper
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&cc); err != nil {
		p.error(w, err)
		return
	}

	config := cc.Config
	hostConfig := cc.HostConfig
	// networkingConfig := cc.NetworkingConfig

	created, err := p.api.ContainerCreate(context.Background(),
		&container.Config{
			Tty:          config.Tty,
			User:         config.User,
			Env:          config.Env,
			Cmd:          config.Cmd,
			AttachStdout: config.AttachStdout,
			AttachStdin:  config.AttachStdin,
			AttachStderr: config.AttachStderr,
			ArgsEscaped:  config.ArgsEscaped,
			Entrypoint:   config.Entrypoint,
			Image:        config.Image,
			WorkingDir:   config.WorkingDir,
		}, &container.HostConfig{
			Privileged:  false,
			AutoRemove:  hostConfig.AutoRemove,
			Cgroup:      container.CgroupSpec(p.cgroup),
			NetworkMode: "default",
		}, &networktypes.NetworkingConfig{
			EndpointsConfig: map[string]*networktypes.EndpointSettings{
				"default": {},
			},
		}, name)
	if err != nil {
		p.error(w, err)
		return
	}

	httputils.WriteJSON(w, http.StatusCreated, created)
}
