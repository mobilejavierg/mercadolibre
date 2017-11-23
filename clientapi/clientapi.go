// Autor: Javier Gonzalez
// Fecha 20-11-2017
// email: mobile.javierg@gmail.como

package clientApi

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"

	//requisitos de googlecloud/appengine
	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/urlfetch"
)

func AsyncGetArticles(wg *sync.WaitGroup, categoriesId string, datos chan Search, offset int, sortId string, req *http.Request) {

	//por un tema de performance solo tomo 100 articulo por Request (como maximo el api devuelve 200)
	_url := fmt.Sprintf("https://api.mercadolibre.com/sites/MLA/search?category=%s&limit=100&offset=%d&sort=%s", categoriesId, offset, sortId)
	///////////////////////////////////////////////////////
	//appEgine adapter

	///	new parameter: req * http.Request

	ctx := appengine.NewContext(req)
	ctx, cancel := context.WithTimeout(ctx, 20*time.Second)

	client := urlfetch.Client(ctx)
	response, err := client.Get(_url)
	/////////////////////////////////////////////////
	//old line	response, err := http.Get(_url)

	if err != nil {
		log.Printf(err.Error())
		return
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

	//envio al channel los datos recuperados
	datos <- tmpDatos

	//Wait Group, herramienta que permite sincronizar los gorouting
	if wg != nil {
		defer func() {
			wg.Done()
			cancel()
		}()
	}

}

func GetPopulation(categoriesId string, req *http.Request) int {
	//objetivo obtener datos generales, mas especificamente la cantidad total de articulos por categoria
	//

	_url := fmt.Sprintf("https://api.mercadolibre.com/sites/MLA/search?category=%s&sort=price_desc&limit=1", categoriesId)

	///////////////////////////////////////////////////////
	//appEgine adapter
	///	new parameter: req * http.Request
	var cancel context.CancelFunc
	var response *http.Response
	var err error

	//con "Go test" no tengo el contexto de appengine
	//debo utilizar http.get
	if req != nil {
		ctx := appengine.NewContext(req)
		ctx, cancel = context.WithTimeout(ctx, 20*time.Second)
		client := urlfetch.Client(ctx)
		response, err = client.Get(_url)
	} else {
		//si req es nil, utilizo el http.get, ya que es llamado por go test
		response, err = http.Get(_url)
	}
	//old line:	response, err := http.Get(_url)
	/////////////////////////////////////////////////

	if err != nil {
		log.Printf(err.Error())
		return 0
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	var tmpDatos Search
	json.Unmarshal(responseData, &tmpDatos)
	totalDatos := tmpDatos.Paginado.Total

	defer func() {
		//si req es distinto de nil
		//llamo a cancel() ya que estoy usando el contexto de appengine
		if req != nil {
			cancel()
		}
	}()

	return totalDatos

}

func GetSampleLen(population int) int {

	/* Formula con variable cuantitativa en poblacion conocida
	   Utilizando como constantes:
	   Z = 95% => 1.96 nivel de confianza
	   e = 5% => 0.05 error estimado
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
	const e float32 = 0.05
	const DevStan float32 = 0.25

	N = float32(population)

	n := (N * (Z * Z) * (DevStan * DevStan)) / ((N - 1) * (e * e) * (Z * Z) * (DevStan * DevStan))

	return int(n)

}

func Analize_data(categoriesId string, req *http.Request) (_respuesta ResponseAPI) {

	var wg sync.WaitGroup
	var done bool
	var globalDatos Search
	arDatos1 := make(chan Search)

	//total de datos
	total := GetPopulation(categoriesId, req)
	//calcular el tamaño de la muestra
	total = GetSampleLen(total)
	//dividir por el offset maximo de 200
	total = total / 200
	done = false

	///////////////////////////////////////////////////////////////////////
	//con este GET obtengo el muestreo con los precios mas altos
	order := "price_desc"
	var offsetAcum int = 0
	wg.Add(1)
	go AsyncGetArticles(&wg, categoriesId, arDatos1, offsetAcum, order, req)

	///////////////////////////////////////////////////////////////////////
	//loop para procesar el resto de las muestras
	//inicio con price_asc para tener los precios minimos
	order = "price_asc"
	for i := 0; i <= total; i++ {

		wg.Add(1)
		if i > (total / 2) {
			order = "price_desc"
		}

		go AsyncGetArticles(&wg, categoriesId, arDatos1, offsetAcum, order, req)
		offsetAcum += 200

	}

	//inicio gorouting, funcion asincronica que permite monitorear los waitgroup y cortar el loop infinito
	go monitorDonde(&wg, &done)

	for !done {

		tmpDatos := <-arDatos1
		globalDatos.Resultados = append(globalDatos.Resultados, tmpDatos.Resultados...)
		time.Sleep(time.Millisecond * 1)

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

	return
}

func monitorDonde(wg *sync.WaitGroup, done *bool) {
	////////////////////////////////////////////////////////////
	//espero asincronicamente que termine el grupo de goroutines
	wg.Wait()
	*done = true //flag para terminar el loop infinito de analize Analize_data
}

//////////////////////////////////////////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////////////////////////////////////////////////////////////////////////
