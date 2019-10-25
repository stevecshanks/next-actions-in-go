FROM golang:1.13

WORKDIR /go/src/next-actions/api

COPY . .

RUN ["go", "get", "github.com/githubnemo/CompileDaemon"]

ENTRYPOINT CompileDaemon -log-prefix=false -build="go build" -command="./api"