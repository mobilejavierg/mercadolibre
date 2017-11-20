package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/mobilejavierg/mercadolibre/clientapi"
	//	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin"
)

func main() {

	r := gin.Default()
	v1 := r.Group("/categories")
	{
		v1.GET("/:id/prices", getPrices)
	}

	r.Run(":8080")

}

func getPrices(c *gin.Context) {

	id := c.Params.ByName("id")
	resp := analize_data(id)
	c.JSON(200, resp)

}

func analize_data(categoriesId string) (_respuesta clientApi.ResponseAPI) {

	var wg sync.WaitGroup
	var done bool
	var globalDatos clientApi.Search
	arDatos1 := make(chan clientApi.Search)

	//total de datos
	total := clientApi.GetPopulation(categoriesId)
	//calcular el tamaño del muestro
	total = clientApi.GetSampleLen(total)
	//dividir por el offset maximo de 200
	total = total / 200
	done = false

	//con este GET obtengo el muestreo con los precios mas altos
	order := "price_desc"
	var offsetAcum int = 0
	wg.Add(1)
	go clientApi.AsyncGetArticles(&wg, categoriesId, arDatos1, offsetAcum, order)
	//****************************

	//empiezo e loop para procesar el resto de las muestras
	//inicio con price_asc para tener los precios minimos
	order = "price_asc"
	for i := 0; i <= total; i++ {

		wg.Add(1)
		if i > (total / 2) {
			order = "price_desc"
		}

		go clientApi.AsyncGetArticles(&wg, categoriesId, arDatos1, offsetAcum, order)
		offsetAcum += 200

		fmt.Println(i, " de ", total)

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
	_max, _min, _mediana := clientApi.GetEstadistics(globalDatos)

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
