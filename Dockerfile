FROM golang:1.16-alpine
WORKDIR /app
COPY go.mod go.sum main.go ./
RUN go mod download && \
    go build

FROM alpine:latest
WORKDIR /usr/local/bin
COPY --from=0 /app/nodeportlister ./
CMD ["/usr/local/bin/nodeportlister"]
