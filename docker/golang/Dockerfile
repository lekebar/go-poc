ARG VERSION

FROM bitnami/golang:${VERSION}

USER root

ARG UID
ARG GID

# Add daemon home directory
RUN usermod -u ${UID} daemon
RUN groupmod -g ${GID} daemon