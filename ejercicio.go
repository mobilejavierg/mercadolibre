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

	//clientApi.ProcessListArt()

}

func analize_data(categoriesId string) (_respuesta clientApi.ResponseAPI) {

	var wg sync.WaitGroup
	var done bool
	var offsetAcum int = 0
	var globalDatos clientApi.Search
	arDatos1 := make(chan clientApi.Search)

	total, maxPrice := clientApi.GetPopulation(categoriesId)

	total = clientApi.GetSampleLen(total)
	total = total / 200
	done = false

	order := "price_asc"

	for i := 0; i <= total; i++ {

		wg.Add(1)
		if i > (total / 2) {
			order = "price_desc"
			//fmt.Println(order)
		}

		go clientApi.AsyncGetArticles(&wg, "MLA1430", arDatos1, offsetAcum, order)
		offsetAcum += 200

	}
	go monitorDonde(&wg, &done)
	time.Sleep(time.Second * 2)

	for !done {

		tmpDatos := <-arDatos1
		globalDatos.Resultados = append(globalDatos.Resultados, tmpDatos.Resultados...)
		fmt.Printf(" max: %.2f \n", maxPrice)
		//		mapB, _ := json.Marshal(globalDatos.Resultados)
		//		fmt.Println(string(mapB))

	}

	//	fmt.Println("fin.")

	//debido a que el API de MELI no garantiza los mismos resultados por cada GET, tengo que validar el maximo

	/*	if maxPrice < globalDatos.Resultados[len(globalDatos.Resultados)-1].Price {
			maxPrice = globalDatos.Resultados[len(globalDatos.Resultados)-1].Price
		}
	*/
	_max, _min, _mediana := clientApi.GetEstadistics(globalDatos)

	_respuesta.Max = maxPrice
	_respuesta.Min = _min
	_respuesta.Seggested = _mediana

	fmt.Printf("maxPrice: %.2f \n", maxPrice)
	fmt.Printf("max: %.2f \n", _max)
	fmt.Printf("min: %.2f \n", _min)
	fmt.Printf("mediana: %.2f \n", _mediana)

	return
	//wg.Wait()

	/*
		r := gin.Default()
		v1 := r.Group("/categories")
		{
			v1.GET("/:id/prices", getPrices)
		}

		r.Run(":8080")*/
}

/*func getPrices(c *gin.Context) {

	//id := c.Params.ByName("id")

	//	resp := clientApi.GetCategories()
	//clientApi.GetArticles("aaa")
	//c.JSON(200, resp)

	clientApi.ProcessListArt()

}*/
func monitorDonde(wg *sync.WaitGroup, done *bool) {

	wg.Wait()

	fmt.Println("fin monitor")
	*done = true
	fmt.Println(*done)

}
