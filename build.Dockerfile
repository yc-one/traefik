FROM traefik-dev:server_1.7
RUN rm -rf /go/src/github.com/containous/traefik
ARG DOCKER_VERSION=18.09.7
ARG DEP_VERSION=0.5.1
WORKDIR /go/src/github.com/containous/traefik
COPY . /go/src/github.com/containous/traefik
