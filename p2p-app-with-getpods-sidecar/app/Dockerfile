FROM golang:1.12 AS builder

ENV CGO_ENABLED=0
WORKDIR $GOPATH/src/msidecar/app
COPY app.go .
RUN go build -o /go/bin/app

FROM scratch
COPY --from=builder /go/bin/app /app

ENTRYPOINT ["/app"]
