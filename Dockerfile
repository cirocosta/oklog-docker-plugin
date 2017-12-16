FROM golang:alpine as builder

ADD ./VERSION /go/src/github.com/cirocosta/oklog-docker-plugin/VERSION
ADD ./main.go /go/src/github.com/cirocosta/oklog-docker-plugin/main.go
ADD ./vendor /go/src/github.com/cirocosta/oklog-docker-plugin/vendor
ADD ./http /go/src/github.com/cirocosta/oklog-docker-plugin/http
ADD ./driver /go/src/github.com/cirocosta/oklog-docker-plugin/driver
ADD ./docker /go/src/github.com/cirocosta/oklog-docker-plugin/docker

WORKDIR /go/src/github.com/cirocosta/oklog-docker-plugin

RUN set -ex && \
  CGO_ENABLED=0 go build \
        -tags netgo -v -a \
        -ldflags "-X main.version=$(cat ./VERSION) -extldflags \"-static\"" && \
  mv ./oklog-docker-plugin /usr/bin/oklog-docker-plugin

FROM alpine
COPY --from=builder /usr/bin/oklog-docker-plugin /usr/local/bin/oklog-docker-plugin

RUN set -x && \
  apk add --update ca-certificates
