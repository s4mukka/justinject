FROM golang:1.21.6-alpine as build

ENV GOPATH /go

COPY build/* /usr/bin/
RUN chmod +x /usr/bin/justinject

EXPOSE 8080

ENTRYPOINT ["/usr/bin/justinject"]
