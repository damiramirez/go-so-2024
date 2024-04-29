**Lineamiento e Implementación**

**Memoria de Instrucciones**


Esta parte de la memoria será la encargada de obtener de los archivos de pseudo código las
instrucciones y de devolverlas a pedido a la CPU.
Al momento de recibir la creación de un proceso, la memoria de instrucciones deberá leer el archivo
de pseudocódigo indicado y generar las estructuras que el grupo considere necesarias para poder
devolver las instrucciones de a 1 a la CPU según ésta se las solicite por medio del Program Counter.
Ante cada petición se deberá esperar un tiempo determinado a modo de retardo en la obtención de
la instrucción, y este tiempo, estará indicado en el archivo de configuración.

**Esquema de memoria**

**Estructuras**

La memoria al trabajar bajo un esquema de paginación simple estará compuesta principalmente por
2 estructuras principales las cuales son:

    ● Un espacio contiguo de memoria (representado por un array de bytes). Este representará el
espacio de usuario de la misma, donde los procesos podrán leer y/o escribir.

    ● Las Tablas de páginas.

Es importante aclarar que cualquier implementación que no tenga todo el espacio de memoria
dedicado a representar el espacio de usuario de manera contigua será motivo de desaprobación
directa, para esto se puede llegar a controlar la implementación a la hora de iniciar la evaluación.

El tamaño de la memoria siempre será un múltiplo del tamaño de página.

**Comunicación con Kernel, CPU e Interfaces de I/O**

# Creación de proceso
Esta petición podrá venir solamente desde el módulo Kernel, y el módulo Memoria deberá crear las
estructuras administrativas necesarias.

$$
Para inicializar la memoria se debe utilizar SIN excepcion la siguiente función
make([]byte, TamMemoria)
$$

**Finalización de proceso**

Esta petición podrá venir solamente desde el módulo Kernel. El módulo Memoria, al ser finalizado un
proceso, debe liberar su espacio de memoria (marcando los frames como libres pero sin
sobreescribir su contenido).

**Acceso a tabla de páginas**
El módulo deberá responder el número de marco correspondiente a la página consultada.

**Ajustar tamaño de un proceso**
Al llegar una solicitud de ajuste de tamaño de proceso (resize) se deberá cambiar el tamaño del
proceso de acuerdo al nuevo tamaño. Se pueden dar 2 opciones:

# Ampliación de un proceso

Se deberá ampliar el tamaño del proceso al final del mismo, pudiendo solicitarse múltiples páginas.
Es posible que en un punto no se puedan solicitar más marcos ya que la memoria se encuentra llena,
por lo que en ese caso se deberá contestar con un error de Out Of Memory.

# Reducción de un proceso

Se reducirá el mismo desde el final, liberando, en caso de ser necesario, las páginas que ya no sean
utilizadas (desde la última hacia la primera).

# Acceso a espacio de usuario
Esta petición puede venir tanto de la CPU como de un Módulo de Interfaz de I/O, es importante
tener en cuenta que las peticiones pueden ocupar más de una página.

El módulo Memoria deberá realizar lo siguiente:

    ● Ante un pedido de lectura, devolver el valor que se encuentra a partir de la dirección física
pedida.

    ● Ante un pedido de escritura, escribir lo indicado a partir de la dirección física pedida. En caso
satisfactorio se responderá un mensaje de ‘OK’.

Cada petición tendrá un tiempo de espera en milisegundos definido por archivo de configuración.

**Logs mínimos y obligatorios**
Creación / destrucción de Tabla de Páginas: “PID: <PID> - Tamaño: <CANTIDAD_PAGINAS>”

Acceso a Tabla de Páginas: “PID: <PID> - Pagina: <PAGINA> - Marco: <MARCO>”

Ampliación de Proceso: “PID: <PID> - Tamaño Actual: <TAMAÑO_ACTUAL> - Tamaño a

Ampliar: <TAMAÑO_A_AMPLIAR>”

Reducción de Proceso: “PID: <PID> - Tamaño Actual: <TAMAÑO_ACTUAL> - Tamaño a

Reducir: <TAMAÑO_A_REDUCIR>”

Acceso a espacio de usuario: “PID: <PID> - Accion: <LEER / ESCRIBIR> - Direccion

fisica: <DIRECCION_FISICA>” - Tamaño <TAMAÑO A LEER / ESCRIBIR>

Archivo de configuración

**Campo            Tipo        Descripción**
port               Numérico    Puerto en el cual se escuchará la conexión de módulo.
memory_size        Numérico    Tamaño expresado en bytes del espacio de usuario de la memoria.
page_size          Numérico    Tamaño de las páginas en bytes.
instructions_path  String      Carpeta donde se encuentran los archivos de pseudocódigo.
delay_response     Numérico    Tiempo en milisegundos que se deberá esperar antes de responder a las 
                               solicitudes de CPU y FS.

Ejemplo de Archivo de Configuración

{
"port": 8002,
"memory_size": 4096,
"page_size": 16,
"instructions_path": "/home/utnso/mappa-pruebas",
"delay_response": 1000
}



**segundo checkpoint**
● Módulo Memoria:
Se encuentra creado y acepta las conexiones.
Es capaz de abrir los archivos de pseudocódigo y envía las instrucciones al CPU

**tercer checkpoint**
● Módulo Memoria:
Se encuentra completamente desarrollada.

