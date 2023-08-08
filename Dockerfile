FROM docker.io/alpine:3.18@sha256:7144f7bab3d4c2648d7e59409f15ec52a18006a128c733fcff20d3a4a54ba44a

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
