FROM golang:1.22 as builder
WORKDIR /app
COPY . .
RUN GOOS=linux CGO_ENABLED=0 go build -ldflags="-w -s" -o server cmd/server.go

FROM scratch
COPY --from=builder /app/.env .
COPY --from=builder /app/server .
CMD ["./server"]
