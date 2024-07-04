#!/bin/bash

# Verificar que se hayan pasado argumentos
if [ -z "$1" ]; then
    echo "Uso: $0 PROCESS_NAME"
    exit 1
fi

KERNEL_URL="http://localhost:8001"


curl -X PUT -H "Content-Type: application/json" -d "{\"path\": \"/home/utnso/tp-2024-1c-sudoers/procesos/$1\"}" ${KERNEL_URL}/process
