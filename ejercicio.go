package main

import (
	"github.com/mobilejavierg/mercadolibre/clientapi"
	//	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin"
	"google.golang.org/appengine"
)

func main() {

	r := gin.Default()
	v1 := r.Group("/categories")
	{
		v1.GET("/:id/prices", getPrices)
	}

	r.Run(":80")
	appengine.Main()
}

func getPrices(c *gin.Context) {

	id := c.Params.ByName("id")
	resp := clientApi.Analize_data(id)
	c.JSON(200, resp)

}
