FROM python:alpine3.16

RUN apk add curl
RUN apk add sudo

RUN apk add --update docker openrc
RUN rc-update add docker boot

RUN apk add --no-cache --virtual builddeps gcc musl-dev && \
  pip --no-cache-dir install aws-sam-cli awscli && \
  apk del builddeps

RUN adduser -s /bin/bash samcli \
  --disabled-password \
  && echo 'samcli ALL=(ALL) NOPASSWD:ALL' >>/etc/sudoers

USER samcli

WORKDIR /app

COPY . ./