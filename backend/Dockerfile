FROM golang:1.25-alpine AS build

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/mapa-do-corre ./cmd/api

FROM alpine:3.21

RUN apk add --no-cache ca-certificates tzdata \
    && addgroup -S app \
    && adduser -S app -G app

WORKDIR /app

COPY --from=build /out/mapa-do-corre /app/mapa-do-corre

USER app
EXPOSE 8080

CMD ["/app/mapa-do-corre"]
