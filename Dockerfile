FROM golang:1.25.7-alpine AS build

ARG VERSION=dev

WORKDIR /app

# Modules layer
COPY go.mod go.sum ./
RUN go mod download

# Build layer
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags "-X main.version=${VERSION}" -o /myapp ./cmd/app

FROM alpine AS run

COPY --from=build /myapp /myapp
COPY --from=build /app/.env /.env

EXPOSE 8080

CMD ["/myapp"]