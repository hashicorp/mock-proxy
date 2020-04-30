FROM golang:1.14.2 as builder

ENV GOBIN=/go/bin
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GO111MODULE=on

WORKDIR /go/src/github.com/hashicorp/mock-proxy

# Improve build hit rate for slow `go mod download` by copying over .mod and
# .sum files first.
COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download

# Currently, all Go files are in these directories, it may be necessary to add
# more directories to this list over time.
COPY ./cmd cmd
COPY ./pkg pkg
RUN go build -o mock-proxy ./cmd/mock-proxy

FROM alpine:3.11
RUN apk add --no-cache ca-certificates squid dumb-init bash git
WORKDIR /

# Initialize Squid SSL db.
RUN /usr/lib/squid/security_file_certgen -c -s /var/lib/ssl_db -M 4MB

COPY --from=builder /go/src/github.com/hashicorp/mock-proxy/mock-proxy .

COPY ./build/package/docker/configs/squid.conf /etc/squid/squid.conf
COPY ./build/package/docker/configs/squid-ssl.conf /etc/squid/squid-ssl.conf
COPY ./build/package/docker/scripts/squid-icap-init.sh .

ENTRYPOINT ["/usr/bin/dumb-init", "--"]
CMD ["/squid-icap-init.sh"]
