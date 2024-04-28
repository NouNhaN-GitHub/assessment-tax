# build stage
FROM golang:1.22-alpine AS build-stage

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY ./ ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /ktax ./

RUN go test -v ./...


# release stage
FROM gcr.io/distroless/base-debian12 AS release-stage

WORKDIR /

COPY --from=build-stage /ktax /ktax

USER nonroot:nonroot

ENTRYPOINT ["/ktax"]