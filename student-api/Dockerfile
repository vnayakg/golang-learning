FROM golang:1.23-bookworm AS base

RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/nonexistent" \
    --shell "/sbin/nologin" \
    --no-create-home \
    --uid 65532 \
    app-user

WORKDIR $GOPATH/src/app/

COPY . .

RUN go mod download \
    go mod verify \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /main .

FROM scratch

COPY --from=base /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

COPY --from=base /usr/share/zoneinfo /usr/share/zoneinfo

COPY --from=base /etc/passwd /etc/passwd

COPY --from=base /main /main

USER app-user

EXPOSE 8080

CMD ["./main"]
