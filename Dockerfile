FROM golang:latest AS builder

RUN mkdir -p /go/src \
    mkdir -p /go/bin \
    mkdir -p /go/pkg

ENV GOPATH=/go
ENV PATH=$GOPATH/bin:$PATH

RUN mkdir -p $GOPATH/src/app
ADD . $GOPATH/src/app

WORKDIR $GOPATH/src/app
RUN go get .
RUN GOOS=linux GOARCH=amd64 go build -o /go/bin/server main.go

FROM golang:latest

COPY --from=builder /go/bin/server /server
RUN chmod +x /server
ENV TINI_VERSION v0.18.0
ADD https://github.com/krallin/tini/releases/download/${TINI_VERSION}/tini /tini
RUN chmod +x /tini


WORKDIR /
ENTRYPOINT ["/tini", "--", "/server"]

