FROM golang:1.20-alpine3.17 as base

RUN apk add --update --no-cache git build-base

WORKDIR /build
COPY go.mod .
COPY go.sum .
RUN go mod download

FROM base as build-amd64

ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64

WORKDIR /build
COPY . .
RUN go build -o ./bin/ ./...

FROM gcr.io/distroless/static as tradeit

EXPOSE 3000

WORKDIR /
COPY --from=build-amd64 --chown=nonroot:nonroot /build/bin/tradeit-do /tradeit-do
ENTRYPOINT ["/tradeit-do"]