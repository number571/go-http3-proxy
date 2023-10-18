package main

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"net/http"

	"github.com/number571/go-http3-proxy/internal/socks5"

	"github.com/quic-go/quic-go"
	"github.com/quic-go/quic-go/http3"
)

const (
	proxyHost = "127.0.0.1:1080"
	proxyURL  = "socks5://" + proxyHost

	remoteHost = "127.0.0.1:8080"
	remoteURL  = "https://" + remoteHost
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
	dialer, err := socks5.NewDialer(proxyURL)
	if err != nil {
		panic(err)
	}

	client := http.Client{
		Transport: &http3.RoundTripper{
			Dial: func(ctx context.Context, addr string, tlsCfg *tls.Config, cfg *quic.Config) (quic.EarlyConnection, error) {
				proxyConn, err := dialer.DialContext(ctx, "udp", remoteHost)
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
			},
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	req, err := http.NewRequest(
		http.MethodPost,
		remoteURL,
		bytes.NewReader([]byte(`hello, world!`)),
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
}
