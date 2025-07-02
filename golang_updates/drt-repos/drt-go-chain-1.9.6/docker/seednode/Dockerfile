FROM golang:1.20.7 AS builder

RUN apt-get update && apt-get install -y
WORKDIR /go/drt-go-chain
COPY . .
RUN go mod tidy
# Seed node
WORKDIR /go/drt-go-chain/cmd/seednode
RUN go build

# ===== SECOND STAGE ======
FROM ubuntu:22.04
COPY --from=builder /go/drt-go-chain/cmd/seednode /go/drt-go-chain/cmd/seednode

WORKDIR /go/drt-go-chain/cmd/seednode/
EXPOSE 10000
ENTRYPOINT ["./seednode"]
