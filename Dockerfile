FROM golang:1.16 AS builder
WORKDIR /go/src/github.com/superorbital/cludo/
COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download  
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /usr/bin/cludo ./cmd/clientCLI

FROM alpine:latest  
RUN apk --no-cache add ca-certificates
RUN mkdir /etc/cludo
WORKDIR /app/
COPY --from=builder /usr/bin/cludo /usr/bin/cludo
ENTRYPOINT ["cludo"]
