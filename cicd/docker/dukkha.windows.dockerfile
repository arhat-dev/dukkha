ARG ARCH=amd64

FROM arhatdev/builder-go:alpine as builder
# TODO: support multiarch build
FROM mcr.microsoft.com/windows/servercore:ltsc2019
ARG APP=dukkha

LABEL org.opencontainers.image.source https://github.com/arhat-dev/dukkha

ENTRYPOINT [ "/dukkha" ]
