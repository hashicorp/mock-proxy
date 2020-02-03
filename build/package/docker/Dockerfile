FROM golang:1.13.7 as builder

ENV GOBIN=/go/bin
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GO111MODULE=on

WORKDIR /go/src/github.com/hashicorp/vcs-mock-proxy

# Improve build hit rate
COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download

COPY ./ .
RUN go build -o vcs-mock-proxy ./cmd/vcs-mock-proxy

FROM alpine:3.9
RUN apk add --no-cache ca-certificates squid dumb-init bash
WORKDIR /

COPY --from=builder /go/src/github.com/hashicorp/vcs-mock-proxy/vcs-mock-proxy .
COPY ./build/package/docker/configs/squid.conf /etc/squid/squid.conf
COPY ./build/package/docker/scripts/squid-icap-init.sh .

ENTRYPOINT ["/usr/bin/dumb-init", "--"]
CMD ["/squid-icap-init.sh"]
