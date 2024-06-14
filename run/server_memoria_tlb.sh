#!/bin/bash

# Verificar que se haya pasado un argumento
if [ -z "$1" ]; then
    echo "Uso: $0 <dev|prod>"
    exit 1
fi

# Asignar el parámetro a la variable ENV
ENV=$1

if [ -z "$2" ]; then
    echo "Uso: $0 PROCESO/MEMORIA_X"
    exit 1
fi

# Ruta de la carpeta a la que quieres cambiar
target_directory="/home/utnso/tp-2024-1c-sudoers"
kernel_config="/home/utnso/tp-2024-1c-sudoers/kernel/config/config_memoria_tlb.json"
cpu_config="/home/utnso/tp-2024-1c-sudoers/cpu/config/config_memoria_tlb.json"
memory_config="/home/utnso/tp-2024-1c-sudoers/memoria/config/config_memoria_tlb.json"
io_config="/home/utnso/tp-2024-1c-sudoers/entradasalida/config/config_memoria_tlb.json"

commands=(
    "cd $target_directory && make kernel ENV=$ENV C=$kernel_config && exec bash"
    "cd $target_directory && make cpu ENV=$ENV C=$cpu_config && exec bash"
    "cd $target_directory && make memoria ENV=$ENV C=$memory_config && exec bash"
)


# Abrir cada comando en una nueva ventana de terminator
for cmd in "${commands[@]}"; do
    terminator -e "$cmd" &
done

# Que espere a que levante kernel y dsp levanta la IO
sleep 1

terminator -e "cd $target_directory && make entradasalida ENV=$ENV N=SLP1 P=$io_config && exec bash"

./run_memoria.sh $2
