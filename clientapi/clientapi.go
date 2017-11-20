package clientApi

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync"
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

	//	go GetArticles(categoriesId, datos, offset)

}

func ProcessListArt() {

	/*
		defer func() {
			contador := <-countGoRoutines
			contador--
			countGoRoutines <- contador
		}()
	*/
}

func AsyncGetArticles(wg *sync.WaitGroup, categoriesId string, datos chan Search, offset int, sortId string) {

	defer wg.Done()
	//price_desc
	//price_asc, price_desc
	_url := fmt.Sprintf("https://api.mercadolibre.com/sites/MLA/search?category=%s&limit=200&offset=%d&sort=%s", categoriesId, offset, sortId)

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

	//	gloalDatos := <-datos

	/*	if gloalDatos == nil {
		gloalDatos = tmpDatos
	}*/

	//	gloalDatos.Resultados = append(gloalDatos.Resultados, tmpDatos.Resultados...)

	if offset > totalDatos {
		return
	}

	//	offset += 200
	//	fmt.Println(len(gloalDatos.Resultados))

	datos <- tmpDatos
	//	go GetArticles(categoriesId, datos, offset)
	//	defer wg.Done()

}

func GetPopulation(categoriesId string) (int, float64) {

	//	_url := fmt.Sprintf("https://api.mercadolibre.com/sites/MLA/search?category=%s&sort=price_asc&limit=1", categoriesId)
	_url := fmt.Sprintf("https://api.mercadolibre.com/sites/MLA/search?category=%s&sort=price_desc&limit=1", categoriesId)

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
	/*	total = totalDatos
		maxPrice = tmpDatos.Resultados[0].Price*/

	/*	mapB, _ := json.Marshal(tmpDatos)
		fmt.Println(string(mapB))*/

	return totalDatos, tmpDatos.Resultados[0].Price

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
