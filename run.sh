#!/bin/bash

# Obtener la ruta del directorio donde se encuentra el script
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" &> /dev/null && pwd)"

# Definir la URL del kernel
KERNEL_URL="http://localhost:8001"

# Lista de archivos de procesos, relativos al script
procesos=(
    "$SCRIPT_DIR/proceso1.txt"
    "$SCRIPT_DIR/proceso2.txt"
    "$SCRIPT_DIR/proceso3.txt"
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