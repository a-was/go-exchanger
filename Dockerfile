FROM golang:1.24-alpine AS builder

WORKDIR /code

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o app


FROM alpine:3

WORKDIR /app
COPY --from=builder /code/app /app/app

CMD ["./app"]
