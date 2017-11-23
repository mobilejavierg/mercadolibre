package main_test

import (
	"testing"

	"github.com/mobilejavierg/mercadolibre/clientapi"
)

func TestMain(t *testing.T) {

	//envio nil ya que no estamos en el contexto de appengine
	if clientApi.GetPopulation("MLA1430", nil) < 0 {
		t.Error("No trae el total de articulos para la categoria MLA1430")
	}

	if clientApi.GetPopulation("MLA1499", nil) < 0 {
		t.Error("No trae el total de articulos para la categoria MLA1499")
	}

}
