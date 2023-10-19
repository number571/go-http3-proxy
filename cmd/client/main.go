package main

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"net/http"
	"time"

	"github.com/wzshiming/socks5"

	"github.com/quic-go/quic-go"
	"github.com/quic-go/quic-go/http3"
)

var (
	// can be overwritten if used docker-mode (init.go)
	proxyHost  = "127.0.0.1:1080"
	remoteHost = "127.0.0.1:8080"
)

// for QUIC protocol (quic-go)
type sConnWrapper struct {
	net.PacketConn
}

func (c *sConnWrapper) SetReadBuffer(bytes int) error {
	socks5udpConn := c.PacketConn.(*socks5.UDPConn)
	udpConn := socks5udpConn.PacketConn.(*net.UDPConn)

	return udpConn.SetReadBuffer(bytes)
}

func (c *sConnWrapper) SetWriteBuffer(bytes int) error {
	socks5udpConn := c.PacketConn.(*socks5.UDPConn)
	udpConn := socks5udpConn.PacketConn.(*net.UDPConn)

	return udpConn.SetWriteBuffer(bytes)
}

func main() {
	client := http.Client{
		Transport: &http3.RoundTripper{
			Dial: proxyDialer("socks5://" + proxyHost),
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	for {
		req, err := http.NewRequest(
			http.MethodPost,
			"https://"+remoteHost,
			bytes.NewReader([]byte(`hello, server!`)),
		)
		if err != nil {
			panic(err)
		}

		resp, err := client.Do(req)
		if err != nil {
			panic(err)
		}

		result, err := io.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}

		fmt.Println(resp.StatusCode, string(result))
		time.Sleep(time.Second)
	}
}

func proxyDialer(proxyURL string) func(ctx context.Context, addr string, tlsCfg *tls.Config, cfg *quic.Config) (quic.EarlyConnection, error) {
	dialer, err := socks5.NewDialer(proxyURL)
	if err != nil {
		panic(err)
	}

	return func(ctx context.Context, addr string, tlsCfg *tls.Config, cfg *quic.Config) (quic.EarlyConnection, error) {
		proxyConn, err := dialer.DialContext(ctx, "udp", addr)
		if err != nil {
			return nil, err
		}

		remoteAddr, err := net.ResolveUDPAddr("udp", addr)
		if err != nil {
			return nil, err
		}

		connWrapper := &sConnWrapper{proxyConn.(net.PacketConn)}
		earlyConn, err := quic.DialEarly(ctx, connWrapper, remoteAddr, tlsCfg, cfg)
		if err != nil {
			return nil, err
		}

		return earlyConn, nil
	}
}
