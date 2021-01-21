FROM golang:latest

WORKDIR /go/src

COPY .env .
COPY app.yaml .
COPY go.mod .
COPY go.sum .

RUN go mod download

RUN go get github.com/githubnemo/CompileDaemon

RUN apt-get update && apt-get install -y ffmpeg

ENTRYPOINT CompileDaemon -build="go build -o portfolio-server ./src" -command="./portfolio-server"
