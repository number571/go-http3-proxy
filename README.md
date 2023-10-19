# go-http3-proxy

Simple example of proxying requests over the HTTP3/QUIC protocol using socks5/udp server. 

### Running

```bash 
## Terminal_1
$ go run ./cmd/server
> Server is listening...

## Terminal_2
$ go run ./cmd/proxy
> Proxy is listening...

## Terminal_3
## This terminal depends on Terminal_1 & Terminal_2
$ go run ./cmd/client
> 200 echo:'hello, server!'
```

## Dependencies

1. Library with QUIC protocol implementation: https://github.com/quic-go/quic-go
2. Socks5 proxy-server with UDP-support: https://github.com/wzshiming/socks5

## Internal

The dependency with the implemented Socks5 proxy-server has a bug in the ReadFrom function in which the client does not receive a response from the server. This bug cannot be solved in any other way except by importing the repository itself. The error is related to the return of the buffer of the wrong size.
