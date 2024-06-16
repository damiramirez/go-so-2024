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
KERNEL_CONFIG="$TARGET_DIR/kernel/config/config_plani_fifo.json"
CPU_CONFIG="$TARGET_DIR/cpu/config/config_plani.json"
MEMORY_CONFIG="$TARGET_DIR/memoria/config/config_plani.json"
IO_CONFIG="$TARGET_DIR/entradasalida/config/config_slp1.json"

commands=(
    "cd $TARGET_DIR && make kernel ENV=$ENV C=$KERNEL_CONFIG && exec bash"
    "cd $TARGET_DIR && make cpu ENV=$ENV C=$CPU_CONFIG && exec bash"
    "cd $TARGET_DIR && make memoria ENV=$ENV C=$MEMORY_CONFIG && exec bash"
)

# Abrir cada comando en una nueva ventana de terminator
for cmd in "${commands[@]}"; do
    terminator -e "$cmd" &
done

# Esperar a que levante kernel y luego levantar la IO
sleep 2

terminator -e "cd $TARGET_DIR && make entradasalida ENV=$ENV N=SLP1 P=$IO_CONFIG && exec bash"

# Ejecutar run_plani.sh
./run_plani.sh