/*
 * SASC (Sistema Automático de detección de Similaridad de Código) analiza todos los archivos de una extensión
 * definida por el usuario (por defecto se usa la extensión "go") desde el directorio de ejecución (inclusive).
 *
 * El análisis consiste en determinar por cada archivo la frecuencia de todos sus caracteres 
 * (vector n-dimensional de características) y calcular la distancia euclidiana entre ellos 
 * usando dicho vector.
 *
 * Luego se pueden generar dos informes en pantalla:
 * - Agrupación de trabajos (grupos) que se encuentran a una distancia máxima definida por el usuario de un código 
 *   central identificado automáticamente.
 *   Disponible desde la versión 1.8. En la versión 2.0 se optimizó en velocidad y en la cantidad de grupos.
 *   Solamente se repiten integrantes en un grupo si aparece un integrante nuevo a una distancia máxima del código central, 
 *   además se marcan los integrantes que están en otros grupos con un (*) y cual es el código que está como centro de 
 *   dicho grupo.
 * - Distancia de cada archivo a todos los demás (entre menor la distancia, mayor la similaridad, en donde 0.0 indica que son idénticos con respecto
 *   a su vector de características).
 *
 * El usuario puede generar un informe en un archivo:
 * - El usuario puede solicitar la generación de un archivo CSV con la matriz (simétrica) de distancias entre los programas.
 *   Esta matriz puede ser visulizada en una hoja electrónica o procesada por algún programa especializado.
 *
 * Autor: Julián Esteban Gutiérrez Posada
 * Fecha: Agosto de 2021
 * Versión: 2.0
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

// Estructura para almacenar la información de la distancia a un archivo.
// Necesario porque al ordenar sin perder la información del código del que se tiene esa distancia
// - indice del código fuente
// - distancia al código fuente con el índice indicado
type Distancia struct {
	indiceCodigoFuente int
	distancia          float64
}

// Estructura para almacenar la información de un archivo
// - nombre del archivo
// - caracteristicas (frecuencias por cada entrada de la tabla ASCII)
// - distancias a todos los demás archivos
// - si el archivo ya pertenece o no a un grupo
type CodigoFuente struct {
	nombre          string
	caracteristica  []int
	tablaDistancias []Distancia
	perteneceGrupo  bool
}

/*
 * Función para obtener los valores por defecto de los parámetros de la aplicación.
 * Por defecto se asume la extensión "go" y sin un valor mínimo de distancia para filtrar la impresión.
 * El usuario puede indicar otra extensión y si lo desea puede definir un valor mínimo
 * return: extensión por defecto, el valor mínimo de la distancia y el nombre del archivo csv
 */
func obtenerValorPorDefecto() (string, float64, string) {
	extensionPorDefecto := "go"
	distanciaMinima := math.MaxFloat64 // Sin distancia máxima
	nombreTablaCSV := ""

	if len(os.Args) >= 2 && len(os.Args) <= 3 {
		extensionPorDefecto = os.Args[1]

		if extensionPorDefecto == "--help" {
			fmt.Println("AYUDA:\n")
			fmt.Println("El programa se puede ejecutar con hasta con dos parámetros opcionales\n")
			fmt.Println("\t ./SASC [extensión] [distancia máxima | nombreTabla.csv]\n")
			fmt.Println("Por defecto se asume \"go\", sin distancia máxima y sin archivo CSV.\n")
			os.Exit(0)
		}

		if len(os.Args) == 3 {
			// Intento de convertir el tercer parámetro a un entero,
			// si es posible, entonces será la distancia máxima definida por el usuario
			// en otro caso será el nombre del archivo CSV
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
 * La lista incluye todos los archivos del directorio actual y todos sus subdirectorio.
 * param: la extensión que deben cumplir para ser ingresados a la lista
 * return: el arreglo con los nombres de todos los archivos que cumplen las condiciones
 */
func obtenerListado(directorioActual string, extension string) ([]string, error) {
	var archivos []string
	var nombre string

	err := filepath.Walk(directorioActual,
		func(path string, info os.FileInfo, err error) error {
			if !info.IsDir() &&
				strings.HasSuffix(path, extension) {
				nombre = strings.Replace(path, directorioActual, ".", 1)
				archivos = append(archivos, nombre)
			}
			return nil
		})

	return archivos, err
}

/*
 * Función para procesar un archivo (determinar la frecuencia de todos los elementos de la tabla ASCII en el archivo)
 * param: nombre del archivo a procesar
 * return: arreglo con la frecuancia de todos los elementos de la tabla ASCII en el archivo indicado
 */
func prodesarArchivo(nombre string) []int {

	tabla := make([]int, MAX_ASCII)

	filebuffer, err := ioutil.ReadFile(nombre)
	if err != nil {
		panic(err)
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
 * Función que calcula la distancia euclidiana entre dos archivos usando el arreglo de frecuencias.
 * param: dos elementos de tipo CodigoFuente
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
 * param: arreglo con los nombres de todos los archivos para determinar sus caracteristicas
 * return: arreglo con las caracteristicas de todos los archivos de la lista
 */
func determinarCaracteristicas(listado []string) []CodigoFuente {
	var tablaCodigoFuente []CodigoFuente

	cantidadArchivo := len(listado)

	for _, archivo := range listado {
		arregloDistancia := make([]Distancia, cantidadArchivo)
		tablaCodigoFuente = append(tablaCodigoFuente, CodigoFuente{nombre: archivo, caracteristica: prodesarArchivo(archivo), tablaDistancias: arregloDistancia, perteneceGrupo: false})
	}

	return tablaCodigoFuente
}

/*
 * Función que determina las distancias entre todos los archivos de la tabla de código fuente
 * Como la matriz de distancias es una matriz simétrica, se optimizó su llenado. 
 * param: arreglo de la información de todos los archivos de código fuente
 * return: completa la información en el arreglo de código fuente con la distancia a todos los demás (matriz de similaridad)
 */
func determinarDistanciasEntreArchivos(tablaCodigoFuente []CodigoFuente) []CodigoFuente {

	var distanciaTemp float64
	var i, j int

	cantidadArchivos := len(tablaCodigoFuente)

	for i = 0; i < cantidadArchivos; i++ {
		for j = 0; j <= i; j++ {
			distanciaTemp = calcularDistancia(tablaCodigoFuente[i], tablaCodigoFuente[j])
			tablaCodigoFuente[i].tablaDistancias[j] = Distancia{indiceCodigoFuente: j, distancia: distanciaTemp}
			tablaCodigoFuente[j].tablaDistancias[i] = Distancia{indiceCodigoFuente: i, distancia: distanciaTemp}
		}
	}

	return tablaCodigoFuente
}

/*
 * Función para imprimir los grupos de trabajo que se encuentran a una distancia máxima.
 * Un programa puede estar en varios grupos, lo que significa que él está a una distancia máxima de varios programas.
 * param: arreglo con la información del código fuente de los archivos y la distancia mínina
 */
func imprimitGrupos(tablaCodigoFuente []CodigoFuente, distanciaMinima float64) {
	var nombre, integrantes string
	var cantidadIntegrantes, cantidadGrupos int
	var imprimirGrupo bool

	fmt.Println("\nGRUPOS CON SUS MIEMBROS A UNA DISTANCIA MÁXIMA DE", distanciaMinima, "RESPECTO AL CÓDIGO CENTRAL\n")

	cantidadGrupos = 1
	for _, archivo := range tablaCodigoFuente {
		integrantes = "GRUPO " + strconv.Itoa(cantidadGrupos) + "\n"
		cantidadIntegrantes = 0
		imprimirGrupo = false
		for _, distanciaArchivo := range archivo.tablaDistancias {
			if distanciaArchivo.distancia <= distanciaMinima {
				if tablaCodigoFuente[distanciaArchivo.indiceCodigoFuente].perteneceGrupo == true {
					nombre = "(*) "
				} else {
					imprimirGrupo = true
					nombre = "    "
					tablaCodigoFuente[distanciaArchivo.indiceCodigoFuente].perteneceGrupo = true
				}
				integrantes += ("\t" + nombre + tablaCodigoFuente[distanciaArchivo.indiceCodigoFuente].nombre)

				if archivo.nombre == tablaCodigoFuente[distanciaArchivo.indiceCodigoFuente].nombre {
					integrantes += " <- Código central"
				}
				integrantes += "\n"
				cantidadIntegrantes++
			}
		}
		if cantidadIntegrantes > 1 && imprimirGrupo == true {
			fmt.Println(integrantes)
			cantidadGrupos++
		}
	}
	fmt.Println()
}

/*
 * Función para imprimir las distancias de cada archivo a todos los demás
 * usando el filtro de la distancia máxima
 * param: arreglo con la información del código fuente de los archivos y la distancia mínina
 */
func imprimirDistancias(tablaCodigoFuente []CodigoFuente, distanciaMinima float64) {
	fmt.Println("\nDISTANCIAS\n")

	for _, archivo := range tablaCodigoFuente {
		// Ordena las distancias de forma ascendente
		sort.Slice(archivo.tablaDistancias, func(j, k int) bool {
			return archivo.tablaDistancias[j].distancia < archivo.tablaDistancias[k].distancia
		})

		fmt.Println(archivo.nombre)

		for _, distanciaArchivo := range archivo.tablaDistancias { // Se recorre toda la matriz para imprimir todas las distancias
			if distanciaArchivo.distancia <= distanciaMinima {
				if archivo.nombre != tablaCodigoFuente[distanciaArchivo.indiceCodigoFuente].nombre {
					fmt.Printf("\t%8.2f %s\n", distanciaArchivo.distancia, tablaCodigoFuente[distanciaArchivo.indiceCodigoFuente].nombre)
				}
			}
		}
		fmt.Println()
	}
}

/*
 * Función para guarda en un archivo CSV las distancias de cada archivo a todos los demás
 * param: arreglo con la información del código fuente de los archivos y nombre CSV
 */
func generarArchivoCSV(tablaCodigoFuente []CodigoFuente, nombreTablaCSV string) {
	var nombre string

	ptrArchivo, err := os.Create(nombreTablaCSV)

	if err != nil {
		panic(err)
	}

	fmt.Fprintf(ptrArchivo, "CÓDIGO FUENTE\t%s", nombre)
	for _, archivo := range tablaCodigoFuente {
		fmt.Fprintf(ptrArchivo, "\t%s", archivo.nombre)
	}
	fmt.Fprintf(ptrArchivo, "\n")

	for _, archivo := range tablaCodigoFuente { // Se genera toda la matriz simétrica, en lugar de generar únicamente la mitad de ella.
		fmt.Fprintf(ptrArchivo, "%s\t", archivo.nombre)

		for _, distanciaArchivo := range archivo.tablaDistancias {
			fmt.Fprintf(ptrArchivo, "\t%8.2f", distanciaArchivo.distancia)
		}
		fmt.Fprintf(ptrArchivo, "\n")
	}
}

/*
 * Función principal
 */
func main() {
	fmt.Println("SISTEMA AUTOMÁTICO DE SIMILARIDAD DE CÓDIGO (SASC)")
	fmt.Println("Julián Esteban Gutiérrez Posada")
	fmt.Println("jugutier@uniquindio.edu.co\n")
	fmt.Println("Versión 2.0 - Licencia GNU - GPL v3")
	fmt.Println("Agosto de 2021\n")

	fmt.Println("Para más información user ./SASC --help\n")

	extensionPorDefecto, distanciaMinima, nombreTablaCSV := obtenerValorPorDefecto()

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

	if nombreTablaCSV != "" {
		fmt.Println("Fase 3 de 3: Generando el archivo \"" + nombreTablaCSV + "\"")
		generarArchivoCSV(tablaCodigoFuente, nombreTablaCSV)
	} else {
		fmt.Println("Fase 3 de 3: Imprimiendo distancia entre archivos de forma creciente...")
		if distanciaMinima < math.MaxFloat64 {
			fmt.Println("             incluye listado de grupos por definir una distancia máxima.")
			imprimitGrupos(tablaCodigoFuente, distanciaMinima)
			fmt.Println(" (*) Este código pertence a otros grupos")
		} else {
			fmt.Println("             NO incluye grupos por no definir una distancia máxima")
		}

		imprimirDistancias(tablaCodigoFuente, distanciaMinima)
	}
}