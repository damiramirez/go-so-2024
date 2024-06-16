#!/bin/bash

# Verificar que se hayan pasado argumentos
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
monitor_config="$TARGET_DIR/entradasalida/config/config_io_monitor.json"
generic_config="$TARGET_DIR/entradasalida/config/config_io_generica.json"
keyboard_config="$TARGET_DIR/entradasalida/config/config_io_teclado.json"
kernel_config="$TARGET_DIR/kernel/config/config_io.json"
cpu_config="$TARGET_DIR/cpu/config/config_io.json"
memory_config="$TARGET_DIR/memoria/config/config_io.json"

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

# Ejecutar los comandos adicionales en nuevas ventanas de terminator
terminator -e "cd $TARGET_DIR && make entradasalida ENV=$ENV N=TECLADO P=$keyboard_config && exec bash"
terminator -e "cd $TARGET_DIR && make entradasalida ENV=$ENV N=MONITOR P=$monitor_config && exec bash"
terminator -e "cd $TARGET_DIR && make entradasalida ENV=$ENV N=GENERICA P=$generic_config && exec bash"

./run_io.sh