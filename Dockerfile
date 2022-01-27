FROM docker.io/golang:1.17 AS builder
WORKDIR /go/src/github.com/superorbital/cludo/
COPY go.mod go.mod
COPY go.sum go.sum
RUN --mount=type=cache,target=/go/pkg/mod go mod download
COPY . .
RUN --mount=type=cache,target=/go/pkg/mod  CGO_ENABLED=0 go build -ldflags "$(go run github.com/ahmetb/govvv -flags -pkg github.com/superorbital/cludo/pkg/build)" -o /usr/bin/cludo ./cmd/cludo

FROM docker.io/alpine:latest
RUN --mount=type=cache,target=/var/cache/apk apk add ca-certificates
RUN mkdir /etc/cludo
WORKDIR /
COPY --from=builder /usr/bin/cludo /usr/bin/cludo
COPY ./docker/entrypoint-cludo.sh /entrypoint.sh
ENTRYPOINT [ "/entrypoint.sh" ]
