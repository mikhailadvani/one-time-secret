FROM golang:1.12.9 as builder

WORKDIR /go/src/github.com/mikhailadvani/one-time-secret
COPY . .

ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64

RUN make build_local_linux

FROM alpine:3.8
WORKDIR /
RUN apk update && \
    apk add ca-certificates && \
    rm -rf /var/cache/apk/* && \
    update-ca-certificates
COPY --from=builder /go/src/github.com/mikhailadvani/one-time-secret/build/one-time-secret-http-linux /one-time-secret
USER nobody
ENTRYPOINT ["/one-time-secret"]
