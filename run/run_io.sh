#!/bin/bash

echo "IO TEST"

# Obtener la ruta del directorio donde se encuentra el script
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" &> /dev/null && pwd)"
PROCESOS_DIR="$SCRIPT_DIR/../procesos"

# Definir la URL del kernel
KERNEL_URL="http://localhost:8001"

# Lista de archivos de procesos, relativos al script
procesos=(
    "$PROCESOS_DIRs/IO_A.txt"
    "$PROCESOS_DIRs/IO_B.txt"
    "$PROCESOS_DIRs/IO_C.txt"
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