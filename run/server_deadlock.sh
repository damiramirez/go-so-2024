#!/bin/bash

# Verificar que se haya pasado un argumento
if [ -z "$1" ]; then
    echo "Uso: $0 <dev|prod>"
    exit 1
fi

# Asignar el parámetro a la variable ENV
ENV=$1

# Definir las rutas de configuración
target_directory="/home/utnso/tp-2024-1c-sudoers"
kernel_config="/home/utnso/tp-2024-1c-sudoers/kernel/config/config_deadlock.json"
cpu_config="/home/utnso/tp-2024-1c-sudoers/cpu/config/config_deadlock.json"
memory_config="/home/utnso/tp-2024-1c-sudoers/memoria/config/config_deadlock.json"
io_config="/home/utnso/tp-2024-1c-sudoers/entradasalida/config/config_espera.json"

# Definir los comandos a ejecutar
commands=(
    "cd $target_directory && make kernel ENV=$ENV C=$kernel_config && exec bash"
    "cd $target_directory && make cpu ENV=$ENV C=$cpu_config && exec bash"
    "cd $target_directory && make memoria ENV=$ENV C=$memory_config && exec bash"
)

# Abrir cada comando en una nueva ventana de terminator
for cmd in "${commands[@]}"; do
    terminator -e "$cmd" &
done

# Esperar un segundo para asegurar que las ventanas anteriores se hayan abierto
sleep 1

# Ejecutar el comando adicional en una nueva ventana de terminator
terminator -e "cd $target_directory && make entradasalida ENV=$ENV N=ESPERA P=$io_config && exec bash"

./run_deadlock.sh