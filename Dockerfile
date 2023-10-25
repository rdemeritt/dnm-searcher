# Build image
FROM gcr.io/mission-e/build/language/go/alpine:latest as builder
ARG APP_SRC=/go/src/code.8labs.io/rdemeritt/dnm-searcher
COPY . ${APP_SRC}
RUN env GOOS=linux GOARCH=amd64 CGO_ENABLED=1 go env; make -C ${APP_SRC} build \
    && ls ${APP_SRC}

# Deployment image
FROM alpine:latest
ARG APP_SRC=/go/src/code.8labs.io/rdemeritt/dnm-searcher
RUN apk --no-cache add ca-certificates
COPY --from=builder ${APP_SRC}/dnm-searcher /usr/local/bin/
ENTRYPOINT ["/bin/sh", "-c", "sleep 15; /usr/local/bin/dnm-searcher"]
