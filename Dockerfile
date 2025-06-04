FROM golang:1.24.3-alpine3.21 as build-env
RUN apk add --no-cache git gcc
RUN mkdir /app
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o openapi
FROM alpine:3.21
COPY --from=build-env /app/openapi .
EXPOSE 8080/tcp
USER 1001
ENTRYPOINT ["./openapi"]