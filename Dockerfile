FROM alpine:latest

ENV MANGAL_DOWNLOADER_PATH=/downloads
ENV MANGAL_USER=abc
ENV MANGAL_UID=1000
ENV MANGAL_GID=1000

WORKDIR "/config"
RUN mkdir -p "${MANGAL_DOWNLOADER_PATH}" && addgroup -g "${MANGAL_GID}" "${MANGAL_USER}" && adduser \
    --disabled-password \
    --gecos "" \
    --home "$(pwd)" \
    --ingroup "${MANGAL_USER}" \
    --no-create-home \
    --uid "${MANGAL_UID}" \
    "${MANGAL_USER}" && \
    chown abc:abc /config "${MANGAL_DOWNLOADER_PATH}"

COPY mangal /usr/local/bin/mangal
RUN chmod +x /usr/local/bin/mangal
USER "${MANGAL_USER}"
ENTRYPOINT ["/usr/local/bin/mangal"]
