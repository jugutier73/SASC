# SASC
El sistema SASC permite determinar la distancia euclidiana entre TODOS los archivos indicados en un grupo de subdirectorios.

Pregunta 
--------

¿Cómo determinar que tan diferentes son los trabajos de programación que un grupo de estudiantes entrega como respuesta a una tarea o un proyecto?

Una posible solución es determinar por cada archivo un conjunto de caracteristicas y usar estas características para determinar la distancia euclidiana que existe entre todos y cada uno de los archivos indicados. Para este caso, el conjunto de caracteristicas se reduce a la frecuencia utilización de todos los caracteres de la tabla ASCII dentro del archivo.

Una distancia de cero entre un par de archivos, indicaría un uso en igual frecuencia de TODOS los caracteres de la tabla ASCII por cada archivo. Sin embargo esto NO NECESARIAMENTE significa que son identicos, pero si es un llamado de atención para un revisión más detalla por parte del docente. Por ejemplo, dos archivos con distancia cero pueden ser uno con la palabra "CASA" y otro con la palabra "SAAC".

Por otro lado, si se analiza código fuente y se detectan dos archivos con la misma frecuencia en todos los caracteres (letras, números, simbolos, tabuladores, espacios, ...) deja algunas inquietudes y merece hacer una verificación manual, en contraste con archivos que son muy distantes (diferentes).


**El sistema SASC es una implementación para resolver esta problemática.**


SISTEMA AUTOMÁTICO DE SIMILARIDAD DE CÓDIGO (SASC)
--------------------------------------------------
SISTEMA AUTOMÁTICO DE SIMILARIDAD DE CÓDIGO (SASC)
Julián Esteban Gutiérrez Posada
jugutier@uniquindio.edu.co

Versión 1.5 - Licencia GNU - GPL v3
Agosto de 2021

Permite determinar la distancia que hay entre múltiples archivos distribuidos en subdirectorios.

Fase 1 de 3: Calculando características de cada archivo...
Fase 2 de 3: Calculando distancia entre los archivos...
Fase 3 de 3: Imprimiendo distancia entre archivos de forma creciente...


Utilización
-----------

- [x] Cree un directorio para almacenar todos los trabajos de los estudiantes, por ejemplo "Proyecto", luego dentro de este directorio cree un directorio por la entrega de cada grupo, por ejemplo "Grupo 01-D" ...

- [x] Ejecute SASC en el directoio base, para el ejemplo "Proyecto" y listo. 
   SASC imprime en la consola un informe con las distancias entre todo par de archivos.

   SASC está programado para buscar por defecto archivo (.go) e imprimir todos los valores de las distancias. Sin embargo, ambos datos pueden ser configurados, por ejemplo:

    a. Analiza en todos los subdirectorios del directorio actual por programas de extensión (.go) e 
       imprimir todas las distancias.

       ./SASC 

    b. Analiza en todos los subdirectorios del directorio actual por programas de extensión (.java) e 
       imprimir todas las distancias.

       ./SASC java

    b. Analiza en todos los subdirectorios del directorio actual por programas de extensión (.c) e 
       imprimir solamente las distancia menores o iguales 30, recuerde que 0 es que ambos archivos usaron la misma cantidad todos los caracteres de la tabla ASCII

       ./SASC c 30
       
    c. Analiza en todos los subdirectorios del directorio actual por programas de extensión (.go) y 
       crea un archivo llamada reporte.csv con una tabla con las distancias en formato CSV.

       ./SASC go repore.csv


Si lo desea, puede redireccionar la salida de SASC a un archivo de texto a nivel de informe, por ejemplo:

    ./SASC java 40 >informe.txt

SASC está compilado en versión para macOS, Windows 64 y Linux de 64 bits


EJEMPLO
-------

1) Suponga que tres grupo entregaron los archivos: sumaFlujos.go, sumaTubos.go, sumarF.go, respectivamente. Estos archivos se guardan en las en los subdirectorios: *"Grupo 01"*, *"Grupo 02"* y *"Grupo 03"* dentro del directorio *"/Trabajo"*. 

Ubicado en el directorio "/Trabajo" se ejecuta el comando:

`./SASC go 100`

Obteniendo la siguiente salida:

`
SISTEMA AUTOMÁTICO DE SIMILARIDAD DE CÓDIGO (SASC)
Julián Esteban Gutiérrez Posada
jugutier@uniquindio.edu.co

Versión 1.5 - Licencia GNU - GPL v3
Agosto de 2021

Para más información user ./SASC --help

Procesando 3 archivo de extensión .go en /Trabajo

Fase 1 de 3: Calculando características de cada archivo...
Fase 2 de 3: Calculando distancia entre los archivos...
Fase 3 de 3: Imprimiendo distancia entre archivos de forma creciente...

./Grupo 01/sumaFlujos.go
	   75.28 ./Grupo 02/sumaTubos.go
	   99.94 ./Grupo 03/sumarF.go

./Grupo 02/sumaTubos.go
	   75.28 ./Grupo 01/sumaFlujos.go
	   83.42 ./Grupo 03/sumarF.go

./Grupo 03/sumarF.go
	   83.42 ./Grupo 02/sumaTubos.go
	   99.94 ./Grupo 01/sumaFlujos.go
`

En este ejemplo hipotético, sería interesante dar una mirada con detenimiento a los archivos
./Grupo 01/sumaFlujos.go y ./Grupo 03/sumarF.go


2) Ubicado en el directorio "/Trabajo" se ejecuta el comando:

`./SASC go Reporte.csv

Procesando 3 archivo de extensión .go en /Trabajo

Fase 1 de 3: Calculando características de cada archivo...
Fase 2 de 3: Calculando distancia entre los archivos...
Fase 3 de 3: Generando el archivo "Reporte.csv"

-> Reporte.csv <-

CÓDIGO FUENTE       	./Grupo 01/sumaFlujos.go    ./Grupo 02/sumaTubos.go    ./Grupo 03/sumarF.go
./Grupo 01/sumaFlujos.go      	 0.00          	        75.28             	99.94
./Grupo 02/sumaTubos.go       	75.28           	 0.00             	83.42
./Grupo 03/sumarF.go          	99.94                 	83.42              	 0.00      
`


Espero les sea de utilidad.


