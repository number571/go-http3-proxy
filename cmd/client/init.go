package main

import (
	"io"
	"net/http"
	"os"
)

// docker use
func init() {
	switch {
	case len(os.Args) != 2:
		return
	case os.Args[1] != "docker":
		return
	}

	// got IP address of server by domain from 'server-proxy' network
	resp, err := http.Get("http://proxy:1090?destination=server:8080")
	if err != nil {
		panic(err)
	}

	host, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	proxyHost = "proxy:1080"
	remoteHost = string(host)
}
