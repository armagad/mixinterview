FROM ubuntu:14.04

ENV DEBIAN_FRONTEND noninteractive

ADD main /usr/local/bin/

ENV DATABASE_URL postgres://172.17.0.1/?user=mix&sslmode=disable


# RUN apt-get -y --no-install-recommends install git #adduser

# ADD golang.linux.tgz /usr/local/
#
# RUN mkdir -p /go/bin
# RUN mkdir -p /go/pkg
# RUN mkdir -p /go/src
#
# ENV GOROOT      /usr/local/go
# ENV PATH        $GOROOT/bin:/go/bin:$PATH
# ENV GOPATH      /go
# ENV CGO_ENABLED 0

# RUN go get github.com/dghubble/go-twitter
# RUN go get github.com/jackc/pgx
# RUN go get github.com/lib/pq
