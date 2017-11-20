package clientApi

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

func GetCategories() []Categories {

	response, err := http.Get("https://api.mercadolibre.com/sites/MLA/categories")

	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	var categoriesObject []Categories

	json.Unmarshal(responseData, &categoriesObject)

	for i := 0; i < len(categoriesObject); i++ {
		fmt.Println("id: " + categoriesObject[i].Id)
		fmt.Println("nombre: " + categoriesObject[i].Name)
	}

	return categoriesObject

}

//https://api.mercadolibre.com/sites/MLA/search?category=MLA5726&sort=price_asc
func GetArticles(categoriesId string, datos *Search, offset int) {

	_url := fmt.Sprintf("https://api.mercadolibre.com/sites/MLA/search?category=%s&sort=price_asc&limit=200&offset=%d", categoriesId, offset)

	fmt.Println(_url)

	response, err := http.Get(_url)
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	var tmpDatos Search

	json.Unmarshal(responseData, &tmpDatos)
	totalDatos := tmpDatos.Paginado.Total

	if datos == nil {
		*datos = tmpDatos
	}

	datos.Resultados = append(datos.Resultados, tmpDatos.Resultados...)

	if offset > totalDatos {
		return
	}

	offset += 200
	fmt.Println(len(datos.Resultados))

}

func AsyncGetArticles(wg *sync.WaitGroup, categoriesId string, datos chan Search, offset int, sortId string) {

	if wg != nil {
		defer wg.Done()
	}
	//price_asc, price_desc
	_url := fmt.Sprintf("https://api.mercadolibre.com/sites/MLA/search?category=%s&limit=200&offset=%d&sort=%s", categoriesId, offset, sortId)

	//fmt.Println(_url)
	response, err := http.Get(_url)
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	var tmpDatos Search

	json.Unmarshal(responseData, &tmpDatos)
	totalDatos := tmpDatos.Paginado.Total

	if offset > totalDatos {
		return
	}

	datos <- tmpDatos

}

func GetPopulation(categoriesId string) int {

	//	_url := fmt.Sprintf("https://api.mercadolibre.com/sites/MLA/search?category=%s&sort=price_asc&limit=1", categoriesId)
	_url := fmt.Sprintf("https://api.mercadolibre.com/sites/MLA/search?category=%s&sort=price_desc&limit=1", categoriesId)

	//	fmt.Println(_url)

	response, err := http.Get(_url)
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	var tmpDatos Search

	json.Unmarshal(responseData, &tmpDatos)

	totalDatos := tmpDatos.Paginado.Total
	return totalDatos

}

func GetSampleLen(population int) int {

	/* Formula con variable cuantitativa en poblacion conocida
	   Utilizando como constantes:
	   Z = 95% => 1.96 nivel de confianza
	   e = 3% => 0.03 error estimado
	   DevStan = 0.25 desviacion standard, valor utilizado ya que desconozco ese dato, por lo tanto se toma una postura conservadora
	   N = tamaño total de la poblacion a estudiar
	   n = tamaño de la muestra a estudiar

	   n =     N . Z^2 . DevStan^2
	   	   =============================
		   (N-1) . e^2 . Z^2 . DevStan^2

	  N = population

	*/
	var N float32
	const Z float32 = 1.96
	const e float32 = 0.01
	const DevStan float32 = 0.25

	N = float32(population)

	n := (N * (Z * Z) * (DevStan * DevStan)) / ((N - 1) * (e * e) * (Z * Z) * (DevStan * DevStan))

	return int(n)

}
func Analize_data(categoriesId string) (_respuesta ResponseAPI) {

	var wg sync.WaitGroup
	var done bool
	var globalDatos Search
	arDatos1 := make(chan Search)

	//total de datos
	total := GetPopulation(categoriesId)
	//calcular el tamaño del muestro
	total = GetSampleLen(total)
	//dividir por el offset maximo de 200
	total = total / 200
	done = false

	//con este GET obtengo el muestreo con los precios mas altos
	order := "price_desc"
	var offsetAcum int = 0
	wg.Add(1)
	go AsyncGetArticles(&wg, categoriesId, arDatos1, offsetAcum, order)
	//****************************

	//empiezo e loop para procesar el resto de las muestras
	//inicio con price_asc para tener los precios minimos
	order = "price_asc"
	for i := 0; i <= total; i++ {

		wg.Add(1)
		if i > (total / 2) {
			order = "price_desc"
		}

		go AsyncGetArticles(&wg, categoriesId, arDatos1, offsetAcum, order)
		offsetAcum += 200

		//fmt.Println(i, " de ", total)

	}
	go monitorDonde(&wg, &done)

	for !done {

		tmpDatos := <-arDatos1
		globalDatos.Resultados = append(globalDatos.Resultados, tmpDatos.Resultados...)
		time.Sleep(time.Millisecond * 1)
		//		fmt.Printf(" max: %.2f \n", maxPrice)

	}

	//debido a que el API de MELI no garantiza el resultado esperado por cada GET, tengo que validar el maximo
	//ejemplo solicito 1 articulo con sortid: price_desc,el mismo get me devuelve 11111111 o 30
	/*	if maxPrice < globalDatos.Resultados[len(globalDatos.Resultados)-1].Price {
			maxPrice = globalDatos.Resultados[len(globalDatos.Resultados)-1].Price
		}
	*/
	_max, _min, _mediana := GetEstadistics(globalDatos)

	_respuesta.Max = _max
	_respuesta.Min = _min
	_respuesta.Seggested = _mediana

	//	fmt.Printf("maxPrice: %.2f \n", maxPrice)
	fmt.Printf("max: %.2f \n", _max)
	fmt.Printf("min: %.2f \n", _min)
	fmt.Printf("mediana: %.2f \n", _mediana)

	return
}

func monitorDonde(wg *sync.WaitGroup, done *bool) {
	wg.Wait()
	*done = true
}
