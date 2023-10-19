package main

import (
	"fmt"

	"github.com/number571/go-http3-proxy/internal/socks5"
)

func main() {
	// TCP is used only to create a `client-proxy` connection,
	// but all other traffic (client-proxy-server) already uses the UDP protocol
	fmt.Println("Proxy is listening...")
	fmt.Println(socks5.NewServer().ListenAndServe("tcp", "0.0.0.0:1080"))
}
