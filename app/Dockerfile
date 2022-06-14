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

RUN go build -a -mod readonly -o gateway cmd/gateway/main.go && \
    go build -a -mod readonly -o auth cmd/auth/main.go && \
    go build -a -mod readonly -o pipeline cmd/pipeline/main.go && \
    chmod +x gateway auth pipeline

FROM gcr.io/distroless/base-debian10 as production

COPY --from=build /usr/go/bin/sql-migrate /
COPY --from=build /var/app/migrations /migrations
COPY --from=build /var/app/gateway /
COPY --from=build /var/app/auth /
COPY --from=build /var/app/pipeline /

USER nobody

EXPOSE 8080

FROM production as auth

ENTRYPOINT ["/auth"]

FROM production as pipeline

ENTRYPOINT ["/pipeline"]

FROM production as gateway

ENTRYPOINT ["/gateway"]