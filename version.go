package lancelot

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/docker/docker/api/server/httputils"
)

func (p *Proxy) version(w http.ResponseWriter, r *http.Request) {
	version, err := p.api.ServerVersion(context.Background())
	if err != nil {
		p.error(w, err)
		return
	}

	w.Header().Set("Api-Version", version.APIVersion)
	w.Header().Set("Docker-Experimental", strconv.FormatBool(version.Experimental))
	w.Header().Set("Ostype", version.Os)
	w.Header().Set("Server", fmt.Sprintf("Lancelot/%s (%s)", version, version.Os))

	httputils.WriteJSON(w, http.StatusOK, version)
}
