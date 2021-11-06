FROM golang:alpine

COPY . /app
WORKDIR /app
COPY .env.example .env

RUN go build -o main cmd/web/main.go

EXPOSE 8888

ENTRYPOINT [ "/app/main" ]
