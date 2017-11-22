package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mobilejavierg/mercadolibre/clientapi"
	"google.golang.org/appengine"
)

func init() {

	r := gin.New()
	v1 := r.Group("/categories")
	{
		v1.GET("/:id/prices", getPrices)
	}

	//tengo que usar http.Handle para acoplarme al appengine
	http.Handle("/", r)

	//r.run genera conflictos con appengine
	//r.Run(":80")
	//appengine.Main()
}

func main() {
	appengine.Main()
}

func getPrices(c *gin.Context) {

	id := c.Params.ByName("id")

	//debo enviar el objeto Request por requerimientos del appengine, al realizar GET's
	resp := clientApi.Analize_data(id, c.Request)
	c.JSON(200, resp)

}
