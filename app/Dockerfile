ARG GO_VERSION=1.18.1

FROM golang:${GO_VERSION} as base

ENV APP_PATH="/var/app" \
    GO111MODULE="on" \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOFLAGS="-mod=vendor" \
    GOPATH="/usr/go"

COPY app/ ${APP_PATH}

WORKDIR ${APP_PATH}

FROM base as build

RUN go install -mod=readonly github.com/rubenv/sql-migrate/...@latest

RUN go build -a -mod readonly -o app cmd/main.go && \
    chmod +x app

FROM gcr.io/distroless/base-debian10 as production

COPY --from=build /usr/go/bin/sql-migrate /
COPY --from=build /var/app/migrations /migrations
COPY --from=build /var/app/app /

USER nobody

EXPOSE 8080
ENTRYPOINT ["/app"]
