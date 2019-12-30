FROM golang:latest

WORKDIR /go/src

COPY . .

RUN go mod download

RUN go get github.com/githubnemo/CompileDaemon

ENTRYPOINT CompileDaemon -build="go build main.go" -command="./main"
