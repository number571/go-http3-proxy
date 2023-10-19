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
> 200 echo:'hello, server!'
> 200 echo:'hello, server!'
...
```

## Dependencies

1. Library with QUIC protocol implementation: https://github.com/quic-go/quic-go
2. Socks5 proxy-server with UDP-support: https://github.com/wzshiming/socks5

## Docker

This example can also be run using docker. In this case, it is enough to use the make command, after which docker-compose will create three services: `server`, `proxy` and `client`. The client and server do not communicate directly with each other, but use bridges: `client-proxy` and `server-proxy`.

### Running

```bash 
$ make
> go-http3-proxy-server-1  | Server is listening...
> go-http3-proxy-proxy-1   | Proxy is listening...
> go-http3-proxy-client-1  | 200 echo:'hello, server!'
> go-http3-proxy-client-1  | 200 echo:'hello, server!'
> go-http3-proxy-client-1  | 200 echo:'hello, server!'
...
```

### Docker-Compose

```yaml
version: "3"
services:
  server:
    build:
      context: ./
      dockerfile: cmd/server/Dockerfile
    networks:
      - server-proxy
  proxy:
    build:
      context: ./
      dockerfile: cmd/proxy/Dockerfile
    networks:
      - server-proxy
      - client-proxy
  client:
    build:
      context: ./
      dockerfile: cmd/client/Dockerfile
    depends_on:
      - proxy
      - server
    networks:
      - client-proxy
networks:
  client-proxy:
    driver: bridge
  server-proxy:
    driver: bridge
```
