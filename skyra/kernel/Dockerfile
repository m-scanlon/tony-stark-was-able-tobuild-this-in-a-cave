FROM golang:1.24-alpine AS build

WORKDIR /src
COPY ./skyra/kernel /src/skyra/kernel
COPY ./skyra/chain-of-thought /src/skyra/chain-of-thought

WORKDIR /src/skyra/kernel
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/skyra-kernel ./cmd

FROM alpine:3.21

RUN addgroup -S skyra && adduser -S skyra -G skyra

WORKDIR /app
COPY --from=build /out/skyra-kernel /usr/local/bin/skyra-kernel
COPY --from=build /src/skyra/chain-of-thought /app/skyra/chain-of-thought

USER skyra

EXPOSE 9090

ENTRYPOINT ["/usr/local/bin/skyra-kernel"]
