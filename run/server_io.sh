#!/bin/bash

# Verificar que se haya pasado un argumento
if [ -z "$1" ]; then
    echo "Uso: $0 <dev|prod>"
    exit 1
fi

# Asignar el parámetro a la variable ENV
ENV=$1

# Ruta de la carpeta a la que quieres cambiar
target_directory="/home/utnso/tp-2024-1c-sudoers"
kernel_config="/home/utnso/tp-2024-1c-sudoers/kernel/config/config_io.json"
monitor_config="/home/utnso/tp-2024-1c-sudoers/entradasalida/config/config_io_monitor.json"
generic_config="/home/utnso/tp-2024-1c-sudoers/entradasalida/config/config_io_generica.json"
keyboard_config="/home/utnso/tp-2024-1c-sudoers/entradasalida/config/config_io_teclado.json"

commands=(
    "cd $target_directory && make kernel ENV=$ENV CONFIG=$kernel_config && exec bash"
    "cd $target_directory && make cpu ENV=$ENV && exec bash"
    "cd $target_directory && make memoria ENV=$ENV && exec bash"
)

# Abrir cada comando en una nueva ventana de terminator
for cmd in "${commands[@]}"; do
    terminator -e "$cmd" &
done

# Que espere a que levante kernel y dsp levanta la IO
sleep 1

terminator -e "cd $target_directory && make entradasalida ENV=$ENV N=TECLADO P=$keyboard_config && exec bash"
terminator -e "cd $target_directory && make entradasalida ENV=$ENV N=MONITOR P=$monitor_config && exec bash"
terminator -e "cd $target_directory && make entradasalida ENV=$ENV N=GENERICA P=$generic_config && exec bash"