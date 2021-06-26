#!/bin/sh

case "${MATRIX_KERNEL}" in
linux)
  case "${MATRIX_ARCH}" in
  armv5 | mips64le | mips64lehf)
    printf "cicd/docker/dukkha.linux.%s.dockerfile" "${MATRIX_ARCH%*hf}"
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
