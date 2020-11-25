package lancelot

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"
)

func (p *Proxy) Ping(w http.ResponseWriter, r *http.Request) {
	ping, err := p.api.Ping(context.Background())
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
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
