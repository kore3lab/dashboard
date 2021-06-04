#!/bin/bash

set -e

if [ "$1" = '' ]; then
  while true; do sleep 1000; done
fi

exec "$@"