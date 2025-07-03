FROM golang:1.20.7 AS builder
LABEL maintainer="dharitri"

WORKDIR /drt-go-chain-proxy
COPY . .

# Proxy
WORKDIR /drt-go-chain-proxy/cmd/proxy
RUN go build -ldflags="-X main.appVersion=$(git describe --tags --long --dirty) -X main.commitID=$(git rev-parse HEAD)"

# ===== SECOND STAGE ======
FROM ubuntu:22.04
RUN apt-get update -y && apt-get upgrade -y

COPY --from=builder /drt-go-chain-proxy/cmd/proxy /drt-go-chain-proxy/cmd/proxy

WORKDIR /drt-go-chain-proxy/cmd/proxy/
EXPOSE 8080
ENTRYPOINT ["./proxy"]
