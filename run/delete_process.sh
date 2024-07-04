#!/bin/bash

if [ -z "$1" ]; then
    echo "Uso: $0 PID"
    exit 1
fi

# Define variables
KERNEL_URL="http://localhost:8001"
PROCESSID=$1

# Realiza la solicitud HTTP
curl -X DELETE "${KERNEL_URL}/process/${PROCESSID}"