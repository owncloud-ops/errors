FROM alpine:3.16@sha256:7580ece7963bfa863801466c0a488f11c86f85d9988051a9f9c68cb27f6b7872

LABEL maintainer="ownCloud DevOps <devops@owncloud.com>"
LABEL org.opencontainers.image.authors="ownCloud DevOps <devops@owncloud.com>"
LABEL org.opencontainers.image.title="errors"
LABEL org.opencontainers.image.url="https://github.owncloud.com/owncloud-ops/errors"
LABEL org.opencontainers.image.source="https://github.owncloud.com/owncloud-ops/errors"
LABEL org.opencontainers.image.documentation="https://github.owncloud.com/owncloud-ops/errors"

ADD dist/errors /bin/errors

RUN addgroup -g 1001 -S app && \
    adduser -S -D -H -u 1001 -h /opt/app -s /bin/bash -G app -g app app

RUN apk --update add --no-cache ca-certificates mailcap && \
    mkdir -p /opt/app/data && \
    chown -R app:app /opt/app && \
    chmod 0750 /opt/app/data && \
    rm -rf /var/cache/apk/* && \
    rm -rf /tmp/*

EXPOSE 8080 8081

USER app

WORKDIR /opt/app
HEALTHCHECK --interval=10s --timeout=5s --start-period=2s --retries=5 CMD ["/bin/errors", "health"]
ENTRYPOINT ["/bin/errors"]
CMD ["server"]
