package main

import (
	"fmt"
	"net"
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

	go func() {
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			udpAddr, err := net.ResolveUDPAddr("udp", r.URL.Query().Get("destination"))
			if err != nil {
				panic(err)
			}
			fmt.Fprint(w, udpAddr.String())
		})
		http.ListenAndServe("0.0.0.0:1090", nil)
	}()
}
