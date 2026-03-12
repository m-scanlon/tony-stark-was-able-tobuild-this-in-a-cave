FROM golang:1.24-alpine AS build

WORKDIR /src
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/internal-api ./cmd

FROM alpine:3.21

RUN addgroup -S octos && adduser -S octos -G octos
USER octos

COPY --from=build /out/internal-api /usr/local/bin/internal-api

EXPOSE 9090

ENTRYPOINT ["/usr/local/bin/internal-api"]
