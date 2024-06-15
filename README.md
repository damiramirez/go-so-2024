## Por Terminar

# Para Todos

- Modularizar
- Manejo de Errores
- Testear
- Agregar Log Obligatorios
- Eliminar Funciones Inutiles y Reemplazar Datos Harcodeados(por ejemplo en los requests)

# Kernel

- Implementar listar procesos para aquellos que se encuentren en colas de bloqueados por recursos
- Implementar eliminacion de procesos y asociarlo al planificador de largo plazo
- Implementar Check Interrupt
- Manejo de FS - Crear la struct para manejar las posibles direcciones fisicas

# Memoria

- Desarrollar para STDIN/OUT

# IO

- Desarrollar FS

# CPU :

- Es capaz de resolver las operaciones:COPY_STRING, IO_STDIN_READ, IO_STDOUT_WRITE.
- Implementar TLB FIFO Y LRU
- Desarrollar MMU
- Desarrollar operaciones FS
- Modificar Registros de la CPU y agregar los que faltan

## Checkpoint Tag

Para cada checkpoint de control obligatorio, se debe crear un tag en el
repositorio con el siguiente formato:

```
checkpoint-{número}
```

Donde `{número}` es el número del checkpoint.

Para crear un tag y subirlo al repositorio, podemos utilizar los siguientes
comandos:

```bash
git tag -a checkpoint-{número} -m "Checkpoint {número}"
git push origin checkpoint-{número}
```

Asegúrense de que el código compila y cumple con los requisitos del checkpoint
antes de subir el tag.
