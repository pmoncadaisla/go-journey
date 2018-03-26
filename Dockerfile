FROM frolvlad/alpine-glibc
RUN apk --no-cache add ca-certificates && \
    rm -rf /var/cache/apk/*

COPY build/app /bin/app
CMD ["app"]