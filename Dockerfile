FROM alpine:latest

ENV APP="snatch"

COPY ./build/${APP}_linux_amd64.zip /root/

RUN apk --no-cache --update add \
    unzip \
    ca-certificates

RUN unzip /root/${APP}_linux_amd64.zip \
    && mv ${APP} /usr/bin/ \
    && chmod +x /usr/bin/${APP} \
    && ${APP} -h

RUN apk del unzip
