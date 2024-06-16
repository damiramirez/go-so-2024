#!/bin/bash

# Verificar que se hayan pasado tres argumentos
if [ -z "$1" ] || [ -z "$2" ] || [ -z "$3" ]; then
    echo "Uso: $0 <dev|prod> MEMORIA_X <fifo|lru>"
    exit 1
fi

# Asignar los parámetros a las variables
ENV=$1
PROCESO_MEMORIA=$2
ALGORITMO=$3

# Obtener la ruta del directorio donde se encuentra el script
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" &> /dev/null && pwd)"
TARGET_DIR="$SCRIPT_DIR/.."

# Definir rutas relativas para las configuraciones
kernel_config="$TARGET_DIR/kernel/config/config_memoria_tlb.json"
cpu_config="$TARGET_DIR/cpu/config/config_memoria_tlb_$ALGORITMO.json"
memory_config="$TARGET_DIR/memoria/config/config_memoria_tlb.json"
io_config="$TARGET_DIR/entradasalida/config/config_memoria_tlb.json"

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
terminator -e "cd $TARGET_DIR && make entradasalida ENV=$ENV N=SLP1 P=$io_config && exec bash"

./run_memoria.sh $PROCESO_MEMORIA
