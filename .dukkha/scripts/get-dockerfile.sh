#!/bin/sh

case "${MATRIX_OS}" in
linux)
  case "${MATRIX_ARCH}" in
  armv5 | mips64le)
    printf "cicd/docker/dukkha.linux.%s.dockerfile" "${MATRIX_ARCH}"
    ;;
  *)
    printf "cicd/docker/dukkha.linux.dockerfile"
    ;;
  esac
  ;;
windows)
  printf "cicd/docker/dukkha.windows.dockerfile"
  ;;
esac
