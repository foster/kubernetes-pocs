FROM golang:1.12
ADD busy.go .
COPY --from=lachlanevenson/k8s-kubectl:v1.10.3 /usr/local/bin/kubectl /usr/local/bin/kubectl
ENTRYPOINT ["go", "run", "busy.go"]
