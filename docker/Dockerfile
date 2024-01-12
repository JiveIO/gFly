ARG GO_VERSION="1.21.5"
ARG ALPINE_VERSION="3.19"
FROM golang:${GO_VERSION}-alpine${ALPINE_VERSION}

EXPOSE 80 443
ENTRYPOINT ["/init"]
WORKDIR /home/gfly/app
VOLUME /home/gfly/app

# Runtime distribution filename.
ARG BUILD_DATE
ARG VCS_REF
ARG BUILD_VERSION

ARG hostUID=1000
ARG hostGID=1000
ENV hostUID=$hostUID
ENV hostGID=$hostGID

# Labels.
LABEL com.jivecode.schema-version="1.0" \
    com.jivecode.build-date=$BUILD_DATE \
    com.jivecode.name="jiveio/gfly-base" \
    com.jivecode.description="gFly" \
    com.jivecode.url="http://www.jivecode.com" \
    com.jivecode.vcs-url="https://github.com/jiveio/gfly-base" \
    com.jivecode.vcs-ref=$VCS_REF \
    com.jivecode.vendor="JiveCode" \
    com.jivecode.version=$BUILD_VERSION \
    com.jivecode.docker.cmd="docker run gfly-base"

# Install and use `bash` shell
RUN apk add --no-cache bash \
	curl

# Create a group and user `gfly`
RUN addgroup -g $hostGID gfly &&\
    adduser -S -u $hostUID -G gfly -h /home/gfly -s /sbin/nologin gfly

# Install Go tools
RUN go install -v github.com/go-critic/go-critic/cmd/gocritic@latest
RUN go install golang.org/x/vuln/cmd/govulncheck@latest
RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN go install github.com/cosmtrek/air@latest
RUN curl -sfL https://raw.githubusercontent.com/securego/gosec/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v2.18.2
RUN curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.55.2
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.16.2/migrate.linux-amd64.tar.gz | tar xvz && \
    mv migrate $(go env GOPATH)/bin/migrate

### Install S6 overlay
RUN curl -LSs -o /tmp/s6-overlay-amd64.tar.gz \
    https://github.com/just-containers/s6-overlay/releases/download/v2.2.0.3/s6-overlay-amd64.tar.gz && \
    tar xzf /tmp/s6-overlay-amd64.tar.gz -C / && \
    rm -f /tmp/s6-overlay-amd64.tar.gz

# Copy configuration files
ADD rootfs /

# The folder need when run `air`, `go`, `gocritic` or some else via `gfly` account
ENV GOPATH=/go
ENV GFLY_HOME=/home/gfly

RUN chown -R $hostUID:$hostGID $GOPATH
RUN chown -R $hostUID:$hostGID $GFLY_HOME

