#!/bin/bash

# Definir la URL del kernel
KERNEL_URL="http://localhost:8001"

# Lista de archivos de procesos
procesos=(
    "/home/utnso/tp-2024-1c-sudoers/proceso2.txt"
    "/home/utnso/tp-2024-1c-sudoers/proceso1.txt"
)

# Crear cada proceso usando la API
for proceso in "${procesos[@]}"; do
    echo "Creando proceso desde el archivo $proceso"
    curl -X PUT "$KERNEL_URL/process" -H "Content-Type: application/json" -d "{\"path\":\"$proceso\"}"
    sleep 1
done

# Hacer una petición PUT a /plani después de iniciar todos los procesos
echo "Enviando petición a /plani"
curl -X PUT "$KERNEL_URL/plani"
