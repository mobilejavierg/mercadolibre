package main_test

import (
	"testing"

	"github.com/mobilejavierg/mercadolibre/clientapi"
)

func BenchmarkGetPopulation(b *testing.B) {

	for i := 0; i < b.N; i++ {
		clientApi.GetPopulation("MLA409431", nil)
	}
}

// falta aplicar a la funcion Analize_data
// el uso de http.client cuando el origen no es appengine
// por eso no esta en el benchmark de la aplicacion
/*
func BenchmarkGetData(b *testing.B) {

	for n := 0; n < b.N; n++ {
		clientApi.Analize_data("MLA409431", nil)
	}
}
*/
