FROM golang:1.21 as build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY cmd/ ./cmd
COPY internal/ ./internal
COPY docs/ ./docs

ENV CGO_ENABLED=0
ENV GOOS=linux

RUN go build -o /video-hosting cmd/app/main.go

FROM alpine:latest

WORKDIR /app

COPY configs/main.yaml ./configs/
COPY --from=build /video-hosting .

CMD ["./video-hosting"]

