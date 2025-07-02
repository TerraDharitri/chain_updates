FROM golang:1.20.7 AS builder

RUN apt-get update && apt-get install -y
WORKDIR /go/drt-go-chain
COPY . .
RUN go mod tidy
# Keygenerator node
WORKDIR /go/drt-go-chain/cmd/keygenerator
RUN go build

# ===== SECOND STAGE ======
FROM ubuntu:22.04
COPY --from=builder /go/drt-go-chain/cmd/keygenerator /go/drt-go-chain/cmd/keygenerator

ENTRYPOINT ["/go/drt-go-chain/cmd/keygenerator/keygenerator"]
