FROM golang:alpine AS builder

ENV APP_HOME /code/calligraphy-forum/server
WORKDIR "$APP_HOME"

COPY go.mod ./
COPY go.sum ./
RUN go env -w GOPROXY=https://goproxy.cn,direct \
    && go mod download

COPY . ./
RUN go build

FROM alpine:latest

COPY --from=builder /code/calligraphy-forum/server/calligraphy-forum /app/calligraphy-forum

WORKDIR /app

EXPOSE 8082
CMD ["./calligraphy-forum-go", "--config", "calligraphy-forum-dev.docker.yaml"]