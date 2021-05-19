FROM golang:1.16.4

WORKDIR /namespace-admission-controller

COPY  . .
RUN GOOS=linux GOARCH=amd64 go build -o admission-controller ./validating-ac/main.go

FROM ubuntu:20.04
COPY --from=0 /namespace-admission-controller/admission-controller /usr/local/bin/admission-controller

ENTRYPOINT /usr/local/bin/admission-controller
