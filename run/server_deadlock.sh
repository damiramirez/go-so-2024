#!/bin/bash

# Verificar que se haya pasado un argumento
if [ -z "$1" ]; then
    echo "Uso: $0 <dev|prod>"
    exit 1
fi

# Asignar el parámetro a la variable ENV
ENV=$1

# Obtener la ruta del directorio donde se encuentra el script
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" &> /dev/null && pwd)"
TARGET_DIR="$SCRIPT_DIR/.."

# Definir rutas relativas para las configuraciones
kernel_config="$TARGET_DIR/kernel/config/config_deadlock.json"
cpu_config="$TARGET_DIR/cpu/config/config_deadlock.json"
memory_config="$TARGET_DIR/memoria/config/config_deadlock.json"
io_config="$TARGET_DIR/entradasalida/config/config_espera.json"

# Definir los comandos a ejecutar
commands=(
    "cd $TARGET_DIR && make kernel ENV=$ENV C=$kernel_config && exec bash"
    "cd $TARGET_DIR && make cpu ENV=$ENV C=$cpu_config && exec bash"
    "cd $TARGET_DIR && make memoria ENV=$ENV C=$memory_config && exec bash"
)

# Abrir cada comando en una nueva ventana de terminator
for cmd in "${commands[@]}"; do
    terminator -e "$cmd" &
done

# Esperar un segundo para asegurar que las ventanas anteriores se hayan abierto
sleep 1

# Ejecutar el comando adicional en una nueva ventana de terminator
terminator -e "cd $TARGET_DIR && make entradasalida ENV=$ENV N=ESPERA P=$io_config && exec bash"

./run_deadlock.sh