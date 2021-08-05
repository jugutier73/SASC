/*
 * SASC (Sistema Automático de Similaridad de Código) analiza todos los archivos de una extensión
 * determinada (go por defecto o la que el usuario indique) desde el directorio de ejecución.
 *
 * El análisis consiste en determinar por cada archivo la frecuencia de todos sus caracteres y
 * calcular la distancia euclidiana entre ellos usando dicha información.
 *
 * Luego imprime dos informes:
 * - Agrupación de trabajos (grupos) que se encuentran a una distancia mìnima definida por el usuario. 
 *   Disponible desde la versión 1.8
 * - Distancia de cada archivo a todos los demás (entre menor la distancia, mayor la similaridad, en donde 0.0 indica que son idénticos).
 *
 * El usuario puede:
 * - Restringir la impresión solamente a los que se encuentren a un distancia mímina.
 * - Puede solicitar la generación de un archivo CSV con la matriz (simétrica) de distancias entre los programas. 
 *   Esta matriz puede ser visulizada en una hoja electrónica o procesada por algún programa especializado.
 *
 * Autor: Julián Esteban Gutiérrez Posada
 * Fecha: Agosto de 2021
 * Versión: 1.8
 * Licencia: GNU GPL v3 (https://www.gnu.org/licenses/gpl-3.0.html)
 */

package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

// Constante que indica el tamaño de la tabla ASCII
const MAX_ASCII = 256

// Estructura para almacenar la distancia a un archivo
// - nombre del archivo al que se determinó la distancia
// - distancia euclidiana en el espacio MAX_ASCII a dicho archivo
type DistanciaCodigo struct {
	nombre    string
	distancia float64
}

// Estructura para almacenar la información de un archivo
// - nombre del archivo
// - caracteristicas (frecuencias por cada entrada de la tabla ASCII)
// - distancias a todos los demás archivos
type CodigoFuente struct {
	nombre          string
	caracteristica  []int
	tablaDistancias []DistanciaCodigo
}

/*
 * Función para obtener los valor por defecto de los parámetros de la aplicación.
 * Por defecto se asume la extensión "go" y sin un valor mínimo de distancia para filtrar la impresión.
 * El usuario puede indicar otra extensión y si lo desea puede definir un valor mínimo
 * return: extensión por defecto, el valor mínimo de la distancia y el nombre del archivo csv
 */
func obtenerValorPorDefecto() (string, float64, string) {
	extensionPorDefecto := "go"
	distanciaMinima := math.MaxFloat64 // Sin distancia mínima
	nombreTablaCSV := ""

	if len(os.Args) >= 2 && len(os.Args) <= 3 {
		extensionPorDefecto = os.Args[1]

		if extensionPorDefecto == "--help" {
			fmt.Println("AYUDA:\n")
			fmt.Println("El programa se puede ejecutar con hasta con dos parámetros opcionales\n")
			fmt.Println("\t ./SASC [extensión] [distancia mínima | nombreTabla.csv]\n")
			fmt.Println("Por defecto se asume \"go\", sin distancia mínima y sin archivo CSV.\n")
			os.Exit(0)
		}

		if len(os.Args) == 3 {
			distancia, err := strconv.ParseFloat(os.Args[2], 64)

			if err != nil {
				nombreTablaCSV = os.Args[2]				
			} else {
				distanciaMinima = distancia
			}
		}
	}

	return extensionPorDefecto, distanciaMinima, nombreTablaCSV
}

/*
 * Función para obtener el listado de todos los archivos que cumplan con la extensión definida.
 * Lista incluye todos los archivos del directorio actual y todos sus subdirectorio.
 * param: la extensión que deben cumplir para ser ingresados a la lista
 * return: el arreglo con los nombres de todos los archivos que cumplen las condiciones
 */
func obtenerListado(directorioActual string, extension string) ([]string, error) {
	var archivos []string

	err := filepath.Walk(directorioActual,
		func(path string, info os.FileInfo, err error) error {
			if !info.IsDir() &&
				strings.HasSuffix(path, extension) {
				archivos = append(archivos, path)
			}
			return nil
		})

	return archivos, err
}

/*
 * Función para procesar un archivo (determinar la frecuencias de todos los elementos de la tabla ASCII)
 * param: nombre del archivo a procesar
 * return: arreglo con la frecuancia de todos los elementos de la tabla ASCII en el archivo indicado
 */
func prodesarArchivo(nombre string) []int {

	tabla := make([]int, MAX_ASCII)

	filebuffer, err := ioutil.ReadFile(nombre)
	if err != nil {
		panic ( err )
	}

	inputdata := string(filebuffer)
	data := bufio.NewScanner(strings.NewReader(inputdata))
	data.Split(bufio.ScanRunes)

	for data.Scan() {
		tabla[int(data.Text()[0])]++
	}

	return tabla
}

/*
 * Función que calcula la distacia euclidiana entre dos archivos usando el arreglo de frecuencias.
 * param: dos elementos de tipo CodigoFiente
 * return: el valor de la distancia euclidiana entre estos dos archivos (códigos fuente)
 */
func calcularDistancia(c1 CodigoFuente, c2 CodigoFuente) float64 {

	suma := 0.0
	for i := 0; i < MAX_ASCII; i++ {
		suma += math.Pow((float64)(c1.caracteristica[i]-c2.caracteristica[i]), 2.0)
	}

	return math.Sqrt(suma)
}

/*
 * Función que determina las caracteristicas de todos los archivos indicados
 * param: arreglo con los nombre de todos los archivos para determinar sus caracteristicas
 * return: arreglo con las caracteristicas de todos los archivos de la lista
 */
func determinarCaracteristicas(listado []string) []CodigoFuente {
	var tablaCodigoFuente []CodigoFuente

	for _, archivo := range listado {
		tablaCodigoFuente = append(tablaCodigoFuente, CodigoFuente{nombre: archivo, caracteristica: prodesarArchivo(archivo)})
	}

	return tablaCodigoFuente
}

/*
 * Función que determina las distancias entre todos los archivos de la tabla de código fuente
 * Como la matriz de distancias es una matriz simétrica, se puede optimizar (futura versión)
 * para hace el código más rápido al disminuir la cantidad de veces que se calcula la distancia.
 * Algo similar a los que se hace en la función de los grupos.
 * param: arreglo la información de todos los archivos de código fuente
 * return: completa la información en el arreglo de código fuente con la distancia a todos los demás de forma ascendente
 */
func determinarDistanciasEntreArchivos(tablaCodigoFuente []CodigoFuente) []CodigoFuente {

	var distanciaTemp float64

	for i, archivo1 := range tablaCodigoFuente {
		for _, archivo2 := range tablaCodigoFuente {
			distanciaTemp = calcularDistancia(archivo1, archivo2)
			tablaCodigoFuente[i].tablaDistancias = append(tablaCodigoFuente[i].tablaDistancias, DistanciaCodigo{nombre: archivo2.nombre, distancia: distanciaTemp})
		}
	}

	return tablaCodigoFuente
}

/*
 * Función para imprimir los grupos de trabaja que se encuentran a una distancia mínima.
 * Un programa puede estar en varios grupos, lo que significa que él está a una distancia mìnima de varios programas,
 * pero esos otros programas no están a esa distancia mínima.
 * param: arreglo con la información del código fuente de los archivo, la distancia mínina y directorio actual
 */
 func imprimitGrupos(tablaCodigoFuente []CodigoFuente, distanciaMinima float64, directioActual string) {
	var nombre, integrantes string
	var cantidadIntegrantes, cantidadGrupos int

	fmt.Println("\nGRUPOS\n")

	cantidadGrupos = 1
	for i, archivo := range tablaCodigoFuente {
		integrantes = "GRUPO " + strconv.Itoa(cantidadGrupos) +"\n"
		cantidadIntegrantes = 0
		for _, distanciaArchivo := range archivo.tablaDistancias[i:] { // Se procesa la mitad de la matriz (por encima de la diagonal principal)
			if distanciaArchivo.distancia <= distanciaMinima {
				nombre = strings.Replace(distanciaArchivo.nombre, directioActual, ".", 1)
				integrantes += ( "\t" + nombre + "\n" )
				cantidadIntegrantes++
			}
		}
		if cantidadIntegrantes > 1 {
			fmt.Println (integrantes)
			cantidadGrupos++
		}			
	}
	fmt.Println()
}

/*
 * Función para imprimir las distancias de cada archivo a todos los demás
 * usando el filtro de la distancia mínima
 * param: arreglo con la información del código fuente de los archivo, la distancia mínina y directorio actual
 */
func imprimirDistancias(tablaCodigoFuente []CodigoFuente, distanciaMinima float64, directioActual string) {
	var nombre, nombre2 string

	fmt.Println("\nDISTANCIAS\n")

	for _, archivo := range tablaCodigoFuente {
		// Ordena las distancias de forma ascendente
		sort.Slice(archivo.tablaDistancias, func(j, k int) bool {
			return archivo.tablaDistancias[j].distancia < archivo.tablaDistancias[k].distancia
		})

		nombre = strings.Replace(archivo.nombre, directioActual, ".", 1)
		fmt.Println(nombre)

		for _, distanciaArchivo := range archivo.tablaDistancias { // Se recorre toda la matriz para imprimir todas las distancias 
			if distanciaArchivo.distancia <= distanciaMinima {
				nombre2 = strings.Replace(distanciaArchivo.nombre, directioActual, ".", 1)
				if archivo.nombre != distanciaArchivo.nombre {
					fmt.Printf("\t%8.2f %s\n", distanciaArchivo.distancia, nombre2)
				} 
			}
		}
		fmt.Println()			
	}
}



/*
 * Función para guarda en un archivo CSV las distancias de cada archivo a todos los demás
 * param: arreglo con la información del código fuente de los archivo, directorio actual y nombre CSV
 */
 func generarArchivoCSV(tablaCodigoFuente []CodigoFuente, directioActual string, nombreTablaCSV string) {
	var nombre string

	ptrArchivo, err := os.Create( nombreTablaCSV )

	if err != nil {
		panic( err )
	}

	fmt.Fprintf ( ptrArchivo, "CÓDIGO FUENTE\t%s", nombre )
	for _, distanciaArchivo := range tablaCodigoFuente[0].tablaDistancias {
		nombre = strings.Replace(distanciaArchivo.nombre, directioActual, ".", 1)
		fmt.Fprintf ( ptrArchivo, "\t%s", nombre )
	}
	fmt.Fprintf (ptrArchivo, "\n" )	

	for _, archivo := range tablaCodigoFuente { // Se genera toda la matriz simétrica, en lugar de generar únicamente la mitad de ella.
		nombre = strings.Replace(archivo.nombre, directioActual, ".", 1)

		fmt.Fprintf ( ptrArchivo, "%s\t", nombre )

		for _, distanciaArchivo := range archivo.tablaDistancias {
				fmt.Fprintf(ptrArchivo, "\t%8.2f", distanciaArchivo.distancia )
		}
		fmt.Fprintf (ptrArchivo, "\n" )			
	}
}

/*
 * Función principal
 */
func main() {
	fmt.Println("SISTEMA AUTOMÁTICO DE SIMILARIDAD DE CÓDIGO (SASC)")
	fmt.Println("Julián Esteban Gutiérrez Posada")
	fmt.Println("jugutier@uniquindio.edu.co\n")
	fmt.Println("Versión 1.8 - Licencia GNU - GPL v3")
	fmt.Println("Agosto de 2021\n")

	fmt.Println("Para más información user ./SASC --help\n")

	extensionPorDefecto, distanciaMinima, nombreTablaCSV := obtenerValorPorDefecto()

	//nombreTablaCSV = "salida.csv"

	directorioActual, _ := os.Getwd()

	listado, err := obtenerListado(directorioActual, extensionPorDefecto)

	if err != nil {
		panic("Error al obtener el listado de los programas.")
	}

	fmt.Println("Procesando", len(listado), "archivo de extensión ."+extensionPorDefecto+" en", directorioActual, "\n")

	fmt.Println("Fase 1 de 3: Calculando características de cada archivo...")
	tablaCodigoFuente := determinarCaracteristicas(listado)

	fmt.Println("Fase 2 de 3: Calculando distancia entre los archivos...")
	tablaCodigoFuente = determinarDistanciasEntreArchivos(tablaCodigoFuente)

	if nombreTablaCSV == "" {
		fmt.Println("Fase 3 de 3: Imprimiendo distancia entre archivos de forma creciente...")
		if distanciaMinima < math.MaxFloat64 { 
			imprimitGrupos(tablaCodigoFuente, distanciaMinima, directorioActual)
		}

		imprimirDistancias(tablaCodigoFuente, distanciaMinima, directorioActual)
	} else {
		fmt.Println("Fase 3 de 3: Generando el archivo \""+nombreTablaCSV+"\"")
		generarArchivoCSV(tablaCodigoFuente, directorioActual, nombreTablaCSV)
	}
}
