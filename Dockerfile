# Build stage
FROM golang:1.19-alpine3.17 AS build
WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o go-wallet-watcher cmd/server/main.go

# Final stage
FROM gcr.io/distroless/static-debian10
COPY --from=build /app/go-wallet-watcher /go-wallet-watcher
COPY .env /
USER nonroot:nonroot
EXPOSE 8080
ENTRYPOINT ["/go-wallet-watcher"]
LABEL org.opencontainers.image.source=https://github.com/ronilsonalves/go-wallet-watcher