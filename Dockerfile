FROM golang:latest

WORKDIR /go/src

COPY .env .
COPY app.yaml .
COPY go.mod .
COPY go.sum .
COPY ./src ./src
COPY ./dist ./dist
COPY ./resources ./resources

RUN go mod download

RUN apt-get update && apt-get install -y ffmpeg

RUN go build -o portfolio-server ./src

ENTRYPOINT ["./portfolio-server"]
