package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io"
	"math/big"
	"net/http"

	"github.com/quic-go/quic-go/http3"
)

func main() {
	mux := http.NewServeMux()

	server := http3.Server{
		Addr:      "0.0.0.0:8080",
		Handler:   mux,
		TLSConfig: generateTLSConfig(),
	}

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet, http.MethodPost:
			// pass
		default:
			fmt.Fprint(w, "method is not supported")
			return
		}

		if r.Method == http.MethodGet {
			fmt.Fprint(w, "hello, world!")
			return
		}

		result, err := io.ReadAll(r.Body)
		if err != nil {
			fmt.Fprint(w, "error:", err.Error())
			return
		}

		fmt.Fprintf(w, "echo:'%s'", string(result))
	})

	fmt.Println("Server is listening...")
	fmt.Println(server.ListenAndServe())
}

// Setup a bare-bones TLS config for the server
func generateTLSConfig() *tls.Config {
	key, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		panic(err)
	}

	template := x509.Certificate{SerialNumber: big.NewInt(1)}
	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &key.PublicKey, key)
	if err != nil {
		panic(err)
	}

	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)})
	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER})

	tlsCert, err := tls.X509KeyPair(certPEM, keyPEM)
	if err != nil {
		panic(err)
	}

	return &tls.Config{
		Certificates: []tls.Certificate{tlsCert},
		NextProtos:   []string{"http3-echo-example"},
	}
}
