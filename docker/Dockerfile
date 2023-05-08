FROM golang:1.18-alpine3.17 AS builder

COPY .. /tit-backend

RUN apk add make \
 && cd /tit-backend \
 && make build \
 && make install-migrate

FROM alpine:3.17

COPY --from=builder /tit-backend/.tmp/main /usr/local/bin/tit-backend

COPY --from=builder /tit-backend/migrations /migrations

COPY --from=builder /go/bin/migrate /usr/local/bin/migrate

COPY docker/docker-entrypoint.sh /entrypoint.sh

ENTRYPOINT [ "/entrypoint.sh" ]
