FROM golang:alpine
ENV CLAM_DATA_PATH="/database"
WORKDIR /app
# Image specific layers under this line
RUN apk add --no-cache clamav rsyslog wget clamav-libunrar
COPY docker/clamav/database ${CLAM_DATA_PATH}
COPY docker/clamav/config/ /etc/clamav/
COPY bin/file-uploader .
COPY migrations/ ./migrations
COPY docker/docker-entrypoint.d /etc/docker-entrypoint.d
COPY docker/docker-entrypoint /usr/local/bin/docker-entrypoint
RUN chmod +x /usr/local/bin/docker-entrypoint
ENTRYPOINT ["/usr/local/bin/docker-entrypoint"]

CMD ["/app/file-uploader"]