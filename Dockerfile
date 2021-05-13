ARG ARCH=
FROM ${ARCH}golang:1.16 AS builder
WORKDIR /go/src/github.com/superorbital/cludo/
RUN go get github.com/ahmetb/govvv
COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -ldflags "$(govvv -flags -pkg github.com/superorbital/cludo/pkg/build)" -o /usr/bin/cludo ./cmd/cludo

FROM ${ARCH}alpine:latest
RUN apk --no-cache add ca-certificates
RUN mkdir /etc/cludo
WORKDIR /app/
COPY --from=builder /usr/bin/cludo /usr/bin/cludo
COPY ./docker/entrypoint-cludo.sh /entrypoint.sh
ENTRYPOINT [ "/entrypoint.sh" ]
