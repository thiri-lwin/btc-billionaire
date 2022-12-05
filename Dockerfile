FROM golang:1.18-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o main main.go

FROM alpine:3.10
WORKDIR /app
COPY --from=builder /app/main .
COPY wait-for.sh .
COPY migrations ./migrations
COPY config.env .

EXPOSE 3000
CMD ["/app/main"]