#!/bin/sh

set -eux

# shell evaluation in env renderer is disabled by default
# we should get what we write as is
foo="$(printf "%s" something)"

export foo
export current_host_kernel="{{ env.HOST_KERNEL }}"
export current_host_arch="{{ env.HOST_ARCH }}"
