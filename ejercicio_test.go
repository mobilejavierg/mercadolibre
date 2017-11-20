package main_test

import (
	"testing"

	"github.com/mobilejavierg/mercadolibre/clientapi"
)

func TestMain(t *testing.T) {

	if clientApi.GetPopulation("MLA1430") < 0 {
		t.Error("No trae el total de articulos para la categoria MLA1430")
	}

	if clientApi.GetPopulation("MLA1499") < 0 {
		t.Error("No trae el total de articulos para la categoria MLA1499")
	}

}
