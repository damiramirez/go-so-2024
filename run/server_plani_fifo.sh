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
kernel_config="/home/utnso/tp-2024-1c-sudoers/kernel/config/config_plani_fifo.json"
io_config="/home/utnso/tp-2024-1c-sudoers/entradasalida/config/config_slp1.json"

commands=(
    "cd $target_directory && make kernel ENV=$ENV CONFIG=$kernel_config && exec bash"
    "cd $target_directory && make cpu ENV=$ENV && exec bash"
    "cd $target_directory && make memoria ENV=$ENV && exec bash"
    "cd $target_directory && make entradasalida ENV=$ENV N=SLP1 P=$io_config && exec bash"
)

# Función para abrir terminator y posicionar la ventana
open_and_position_terminal() {
    local cmd="$1"
    local x="$2"
    local y="$3"
    local width="$4"
    local height="$5"
    
    terminator --geometry=${width}x${height}+${x}+${y} -e "bash -c '$cmd'" &
    sleep 1  # Esperar un segundo para que la ventana se abra
}

# Dimensiones de la pantalla y de las ventanas
screen_width=1920
screen_height=1080
terminal_width=$((screen_width / 2))
terminal_height=$((screen_height / 2))

# Abrir y posicionar las terminales
open_and_position_terminal "${commands[0]}" 0 0 $terminal_width $terminal_height
open_and_position_terminal "${commands[1]}" $terminal_width 0 $terminal_width $terminal_height
open_and_position_terminal "${commands[2]}" 0 $terminal_height $terminal_width $terminal_height
open_and_position_terminal "${commands[3]}" $terminal_width $terminal_height $terminal_width $terminal_height
