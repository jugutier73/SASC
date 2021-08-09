# SASC v2.0

El sistema SASC (Sistema Automático de detección de Similaridad de Código) permite determinar la distancia euclidiana entre TODOS los archivos indicados en un grupo de subdirectorios. Esta distancia se emplea para determinar el grado de similaridad entre ellos e identificar posibles grupos de códigos similares que cumplan con una distancia máxima definida por el usuario a un código central, por cada grupo, que es identificado automáticamente por SASC.

Pregunta 
--------

¿Cómo determinar que tan diferentes son los trabajos de programación que un grupo de estudiantes entrega como respuesta a una tarea o un proyecto?

Una posible solución es determinar por cada archivo un conjunto de caracteristicas y usar estas características para determinar la distancia euclidiana que existe entre todos y cada uno de los archivos indicados. Para este caso, el conjunto de caracteristicas se reduce a la frecuencia de la utilización de todos los caracteres de la tabla ASCII dentro del archivo.

Una distancia de cero entre un par de archivos, indicaría un uso en igual frecuencia de TODOS los caracteres de la tabla ASCII por cada archivo. Sin embargo esto NO NECESARIAMENTE significa que son idénticos, pero sí es un llamado de atención para una revisión más detallada por parte del docente. Por ejemplo, dos archivos con distancia cero pueden ser uno con la palabra "CASA" y otro con la palabra "SAAC".

Por otro lado, si se analiza el código fuente y se detectan dos archivos con la misma frecuencia en todos los caracteres (letras, números, simbolos, tabuladores, espacios, ...) deja algunas inquietudes y merece hacer una verificación manual, en contraste con archivos que son muy distantes (diferentes).


**El sistema SASC es una implementación para resolver esta problemática.**


SISTEMA AUTOMÁTICO DE SIMILARIDAD DE CÓDIGO (SASC)
--------------------------------------------------

```
SISTEMA AUTOMÁTICO DE SIMILARIDAD DE CÓDIGO (SASC)
Julián Esteban Gutiérrez Posada
jugutier@uniquindio.edu.co

Versión 2.0 - Licencia GNU - GPL v3
Agosto de 2021

Para más información user ./SASC --help

Procesando 8 archivo de extensión .go en /Users/jugutier/SASC/SASC-src 

Fase 1 de 3: Calculando características de cada archivo...
Fase 2 de 3: Calculando distancia entre los archivos...
```

`Fase 3 de 3: Imprimiendo distancia entre archivos de forma creciente...
             incluye listado de grupos por definir una distancia máxima.`

o puede generar un reporte en CSV

`Fase 3 de 3: Generando el archivo "Reporte.csv"`



Utilización
-----------

- [x] Cree un directorio para almacenar todos los trabajos de los estudiantes, por ejemplo "Proyecto", luego dentro de este directorio cree un directorio por la entrega de cada grupo, por ejemplo "G01-D" ...

- [x] Ejecute SASC en el directoio base, para el ejemplo "Proyecto" y listo. 
   SASC imprime en la consola un informe con las distancias entre todo par de archivos y si definió una distancia máxima, se genera un listado de grupos (disponible desde la versión 2.0)

   SASC está programado para buscar por defecto archivo (.go) e imprimir todos los valores de las distancias. Sin embargo, ambos datos pueden ser configurados, por ejemplo:

    a. Analiza, en todos los subdirectorios y en el directorio actual, todos los programas de extensión (.go) e imprime todas las distancias entre ellos.

       ./SASC 

    b. Analiza, en todos los subdirectorios y en el directorio actual, todos los programas de extensión (.java) e imprime todas las distancias entre ellos.

       ./SASC java

    b. Analiza, en todos los subdirectorios y en el directorio actual, todos los programas de extensión (.c) e imprime solo los programas que están a una distancia menor o igual a 30. 

Recuerde que 0 es que ambos archivos usaron la misma cantidad todos los caracteres de la tabla ASCII. Además busca agrupar los programas según las distancias. Esta agrupación no se imprime cuando no hay una distancia mínima.

       ./SASC c 30
       
   c. Analiza, en todos los subdirectorios y en el directorio actual, todos los programas de extensión (.go) e imprime y crea un archivo llamado reporte.csv con una tabla con las distancias en formato CSV.

       ./SASC go repore.csv


Si lo desea, puede redireccionar la salida de SASC a un archivo de texto a nivel de informe, por ejemplo:

    ./SASC java 40 >informe.txt

SASC está compilado en versión para macOS, Windows 64 y Linux de 64 bits


EJEMPLO
-------

1) Suponga que siete grupo entregaron los archivos: DemoE1.go, DemoE2.go, DemoE3.go, ..., DemoE7.go. Estos archivos se guardan en las en los subdirectorios: *"E1"*, *"E2"* ... *"E7"* dentro del directorio *"/Trabajo"*. En dicho directorio hay un programa llamado SASC.go 

Ubicado en el directorio "/Trabajo" se ejecuta el comando:

`./SASC go 300`

Obteniendo la siguiente salida:

```
SISTEMA AUTOMÁTICO DE SIMILARIDAD DE CÓDIGO (SASC)
Julián Esteban Gutiérrez Posada
jugutier@uniquindio.edu.co

Versión 2.0 - Licencia GNU - GPL v3
Agosto de 2021

Para más información user ./SASC --help

Procesando 8 archivo de extensión .go en /Users/jugutier/SASC/SASC-src 

Fase 1 de 3: Calculando características de cada archivo...
Fase 2 de 3: Calculando distancia entre los archivos...
Fase 3 de 3: Imprimiendo distancia entre archivos de forma creciente...
             incluye listado de grupos por definir una distancia máxima.

GRUPOS CON SUS MIEMBROS A UNA DISTANCIA MÁXIMA DE 300 RESPECTO AL CÓDIGO CENTRAL

GRUPO 1
       ./E1/DemoE1.go <- Código central
       ./E2/DemoE2.go
       ./SASC.go

GRUPO 2
   (*) ./E1/DemoE1.go
   (*) ./E2/DemoE2.go <- Código central
       ./E5/DemoE5.go

GRUPO 3
       ./E3/DemoE3.go <- Código central
       ./E4/DemoE4.go

GRUPO 4
   (*) ./E2/DemoE2.go
   (*) ./E5/DemoE5.go <- Código central
       ./E6/DemoE6.go

GRUPO 5
   (*) ./E5/DemoE5.go
   (*) ./E6/DemoE6.go <- Código central
       ./E7/DemoE7.go


 (*) Este código pertence a otros grupos

DISTANCIAS

./E1/DemoE1.go
     238.92 ./E2/DemoE2.go
     295.56 ./SASC.go

./E2/DemoE2.go
     238.92 ./E1/DemoE1.go
     246.98 ./E5/DemoE5.go

./E3/DemoE3.go
      37.47 ./E4/DemoE4.go

./E4/DemoE4.go
      37.47 ./E3/DemoE3.go

./E5/DemoE5.go
     246.98 ./E2/DemoE2.go
     263.61 ./E6/DemoE6.go

./E6/DemoE6.go
     244.02 ./E7/DemoE7.go
     263.61 ./E5/DemoE5.go

./E7/DemoE7.go
     244.02 ./E6/DemoE6.go

./SASC.go
     295.56 ./E1/DemoE1.go

```


2) Ubicado en el directorio "/Trabajo" se ejecuta el comando:

```./SASC go Reporte.csv

SISTEMA AUTOMÁTICO DE SIMILARIDAD DE CÓDIGO (SASC)
Julián Esteban Gutiérrez Posada
jugutier@uniquindio.edu.co

Versión 2.0 - Licencia GNU - GPL v3
Agosto de 2021

Para más información user ./SASC --help

Procesando 8 archivo de extensión .go en /Users/jugutier/SASC/SASC-src 

Fase 1 de 3: Calculando características de cada archivo...
Fase 2 de 3: Calculando distancia entre los archivos...
```
-> Reporte.csv <-

```
CÓDIGO FUENTE     ./E1/DemoE1.go ./E2/DemoE2.go ./E3/DemoE3.go ./E4/DemoE4.go ./E5/DemoE5.go ./E6/DemoE6.go ./E7/DemoE7.go ./SASC.go
./E1/DemoE1.go        0.00   238.92  2314.30  2339.36   475.49   732.95   973.84   295.56
./E2/DemoE2.go      238.92     0.00  2090.78  2115.12   246.98   503.03   743.17   530.06
./E3/DemoE3.go     2314.30  2090.78     0.00    37.47  1856.99  1601.85  1370.63  2594.01
./E4/DemoE4.go     2339.36  2115.12    37.47     0.00  1880.73  1624.82  1392.35  2619.55
./E5/DemoE5.go      475.49   246.98  1856.99  1880.73     0.00   263.61   504.38   760.42
./E6/DemoE6.go      732.95   503.03  1601.85  1624.82   263.61     0.00   244.02  1018.25
./E7/DemoE7.go      973.84   743.17  1370.63  1392.35   504.38   244.02     0.00  1258.17
./SASC.go        295.56   530.06  2594.01  2619.55   760.42  1018.25  1258.17     0.00
```

Espero les sea de utilidad.


