FROM alpine:3.7
RUN mkdir /app 
COPY ./app-exe /app/ 
COPY . crt.* /app/
COPY . key.* /app/
WORKDIR /app
RUN apk update \
        && apk upgrade \
        && apk add --no-cache \
        ca-certificates \
        && update-ca-certificates 2>/dev/null || true
ARG env
ARG config
ARG version
ENV env=${env}
ENV config=${config}
ENV version=${version}
CMD ["/app/app-exe"]