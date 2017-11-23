// Autor: Javier Gonzalez
// Fecha 20-11-2017
// email: mobile.javierg@gmail.como

package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mobilejavierg/clientapi"
	"google.golang.org/appengine"
)

//funcion init requerida por appengine
func init() {

	r := gin.New()
	v1 := r.Group("/categories")
	{
		v1.GET("/:id/prices", getPrices)
	}

	//tengo que usar http.Handle para acoplarme al appengine
	http.Handle("/", r)

	//r.run genera conflictos con appengine
	//r.Run(":8080")
}

func main() {
	//requerido por appengine
	appengine.Main()
}

//proceso el GET del APIRest
func getPrices(c *gin.Context) {

	//obtengo el ID de la Categoria
	id := c.Params.ByName("id")

	//debo enviar el objeto "c.Request" por requerimientos del appengine, al realizar GET's a una URL externa
	resp := clientApi.Analize_data(id, c.Request)
	c.JSON(200, resp)

}
