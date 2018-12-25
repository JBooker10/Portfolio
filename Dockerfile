FROM golang:1.11.4 as builder


RUN go get -u github.com/golang/dep/cmd/dep

ENV GO111MODULE=off

WORKDIR /go/src/app

ADD Gopkg.toml  Gopkg.toml
ADD Gopkg.lock  Gopkg.lock
ADD go.mod  .
ADD go.sum  .

RUN dep ensure --vendor-only

# add source code
ADD src src

# RUN go get

RUN  CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main src/main.go

FROM alpine:3.7

RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*

WORKDIR /root

COPY --from=builder /go/src/app/main .

# expose port
EXPOSE 8080

# run main.go
CMD ["./main"]

