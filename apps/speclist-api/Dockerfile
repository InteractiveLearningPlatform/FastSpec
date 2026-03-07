FROM golang:1.24-alpine AS build

WORKDIR /src

COPY go.mod ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/speclist-api ./cmd/speclist-api

FROM alpine:3.20

RUN apk add --no-cache ca-certificates

WORKDIR /app

COPY --from=build /out/speclist-api /app/speclist-api

EXPOSE 8080

ENTRYPOINT ["/app/speclist-api"]
