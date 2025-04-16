FROM golang:1.22 AS builder

WORKDIR /app
COPY . .

# Static build flags
RUN go mod tidy && \
    CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -a -installsuffix cgo -o main .

FROM scratch
COPY --from=builder /app/main /main
EXPOSE 8085
CMD ["/main"]
