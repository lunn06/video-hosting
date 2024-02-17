FROM golang:1.21

WORKDIR /app

COPY . ./
RUN go mod download

ENV CGO_ENABLED=0
ENV GOOS=linux

RUN go build -o /video-hosting cmd/app/main.go

EXPOSE 8080

CMD ["/video-hosting"]

