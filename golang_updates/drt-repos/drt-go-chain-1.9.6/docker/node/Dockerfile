FROM golang:1.20.7 AS builder

RUN apt-get update && apt-get upgrade -y
WORKDIR /go/drt-go-chain
COPY . .
RUN go mod tidy
# Dharitri node
WORKDIR /go/drt-go-chain/cmd/node
RUN go build -v -ldflags="-X main.appVersion=$(git describe --tags --long --dirty)"

RUN mkdir -p /lib_amd64 /lib_arm64

RUN cp /go/pkg/mod/github.com/TerraDharitri/$(cat /go/drt-go-chain/go.mod | grep drt-chain-vm-v | sort -n | tail -n -1 | awk -F '/' '{print$3}' | sed 's/ /@/g')/wasmer/libwasmer_linux_amd64.so /lib_amd64/
RUN cp /go/pkg/mod/github.com/TerraDharitri/$(cat /go/drt-go-chain/go.mod | grep drt-go-chain-vm | sort -n | tail -n -1 | awk -F '/' '{print$3}' | sed 's/ /@/g')/wasmer2/libvmexeccapi.so /lib_amd64/

RUN cp /go/pkg/mod/github.com/TerraDharitri/$(cat /go/drt-go-chain/go.mod | grep drt-chain-vm-v | sort -n | tail -n -1 | awk -F '/' '{print$3}' | sed 's/ /@/g')/wasmer/libwasmer_linux_arm64_shim.so /lib_arm64/
RUN cp /go/pkg/mod/github.com/TerraDharitri/$(cat /go/drt-go-chain/go.mod | grep drt-go-chain-vm | sort -n | tail -n -1 | awk -F '/' '{print$3}' | sed 's/ /@/g')/wasmer2/libvmexeccapi_arm.so /lib_arm64/

# ===== SECOND STAGE ======
FROM ubuntu:22.04
ARG TARGETARCH
RUN apt-get update && apt-get upgrade -y
COPY --from=builder "/go/drt-go-chain/cmd/node/node" "/go/drt-go-chain/cmd/node/"

# Copy architecture-specific files
COPY --from=builder "/lib_${TARGETARCH}/*" "/lib/"

WORKDIR /go/drt-go-chain/cmd/node/
EXPOSE 8080
ENTRYPOINT ["/go/drt-go-chain/cmd/node/node"]
