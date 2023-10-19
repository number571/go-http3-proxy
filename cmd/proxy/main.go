package main

import (
	"fmt"

	"github.com/number571/go-http3-proxy/internal/socks5"
)

func main() {
	fmt.Println("Proxy is listening...")
	server := socks5.NewServer()
	err := server.ListenAndServe("tcp", "0.0.0.0:1080")
	if err != nil {
		panic(err)
	}
}
