FROM golang:1.12 AS builder

ENV CGO_ENABLED=0
WORKDIR $GOPATH/src/msidecar/sidecar
COPY sidecar.go .
RUN go get ./... && go build -o /go/bin/sidecar

FROM scratch
COPY --from=builder /go/bin/sidecar /sidecar

ENTRYPOINT ["/sidecar"]
