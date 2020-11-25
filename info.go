package lancelot

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/docker/docker/api/server/httputils"
	moby "github.com/docker/docker/api/types"
)

func (p *Proxy) info(w http.ResponseWriter, r *http.Request) {
	info, err := p.api.Info(context.Background())
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// FIXME w.Header().Set("Api-Version", info.APIVersion)
	// FIXME w.Header().Set("Docker-Experimental", strconv.FormatBool(info.Experimental))
	w.Header().Set("Ostype", info.OSType)
	w.Header().Set("Server", fmt.Sprintf("Lancelot/%s (%s)", version, info.OSType))

	httputils.WriteJSON(w, http.StatusOK,
		// FIXME review which info should be usefull
		moby.Info{
			ID:               info.ID,
			KernelVersion:    info.KernelVersion,
			OperatingSystem:  info.OperatingSystem,
			OSType:           info.OSType,
			Architecture:     info.Architecture,
			ServerVersion:    info.ServerVersion,
			ContainerdCommit: info.ContainerdCommit,
			RuncCommit:       info.RuncCommit,
			InitCommit:       info.InitCommit,
		})
}
