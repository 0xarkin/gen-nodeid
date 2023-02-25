FROM golang:1.18-alpine as builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY main.go ./
RUN go build -o main .

FROM alpine:latest
WORKDIR /
COPY --from=builder /app/main .
CMD ["/main"]