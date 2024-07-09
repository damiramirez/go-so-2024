#!/bin/bash

if [ -z "$1" ]; then
    echo "Uso: $0 PID"
    exit 1
fi

if [ -z "$KERNEL_PORT" ]; then
    echo "The KERNEL_PORT is not set"
    echo "Using default port 8001"
    KERNEL_PORT=8001
fi

if [ -z "$KERNEL_HOST" ]; then
    echo "The KERNEL_HOST is not set"
    echo "Using default host localhost"
    KERNEL_HOST=localhost
fi

KERNEL_URL="http://$KERNEL_HOST:$KERNEL_PORT"
PROCESSID=$1

# Realiza la solicitud HTTP
curl -X DELETE "${KERNEL_URL}/process/${PROCESSID}"