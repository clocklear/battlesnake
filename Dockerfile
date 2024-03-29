FROM golang:1.16-alpine AS build

WORKDIR /go/src
COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -ldflags '-extldflags "-static"' -o battlesnake ./cmd/battlesnake

EXPOSE 8080

FROM scratch

COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=build /go/src/battlesnake /usr/bin/battlesnake

CMD ["battlesnake"]
