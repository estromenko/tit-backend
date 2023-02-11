FROM debian:bullseye

RUN apt-get update && apt-get install -y \
    fluxbox \
    websockify \
    xvfb \
    x11vnc

COPY docker-entrypoint.sh /entrypoint.sh

ENTRYPOINT [ "/entrypoint.sh" ]
