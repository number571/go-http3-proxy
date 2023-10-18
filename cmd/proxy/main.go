// This code copied from: github.com/wzshiming/socks5/cmd/socks5
package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	// origin: github.com/wzshiming/socks5
	"github.com/number571/go-http3-proxy/internal/socks5"
)

var address string
var username string
var password string

func init() {
	flag.StringVar(&address, "a", ":1080", "listen on the address")
	flag.StringVar(&username, "u", "", "username")
	flag.StringVar(&password, "p", "", "password")
	flag.Parse()
}

func main() {
	logger := log.New(os.Stderr, "[socks5] ", log.LstdFlags)
	svc := &socks5.Server{
		Logger: logger,
	}
	if username != "" {
		svc.Authentication = socks5.UserAuth(username, password)
	}
	fmt.Println("Proxy is listening...") // append println
	err := svc.ListenAndServe("tcp", address)
	if err != nil {
		logger.Println(err)
	}
}
