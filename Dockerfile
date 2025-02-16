FROM golang:1.23

ARG CONFIG_PATH

WORKDIR /go/src/avito-shop/

COPY . .

RUN go mod tidy && go build -o /build ./cmd

RUN go clean -cache -modcache

RUN go test ./... -tags=unit

EXPOSE 8080

CMD ["/build"]
