FROM golang:1.12 AS builder

ENV CGO_ENABLED=0
WORKDIR $GOPATH/src/msidecar/sidecar
COPY sidecar.go .
RUN go get ./... && go build -o /go/bin/sidecar

FROM busybox
COPY --from=lachlanevenson/k8s-kubectl:v1.10.3 /usr/local/bin/kubectl /usr/local/bin/kubectl
COPY --from=builder /go/bin/sidecar /sidecar

ENTRYPOINT ["/sidecar"]
