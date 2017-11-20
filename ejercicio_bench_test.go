package main_test

import (
	"testing"

	"github.com/mobilejavierg/mercadolibre/clientapi"
)

func BenchmarkGetPopulation(b *testing.B) {
	for i := 0; i < b.N; i++ {
		clientApi.GetPopulation("MLA409431")
	}
}

func BenchmarkGetData(b *testing.B) {

	for n := 0; n < b.N; n++ {
		clientApi.Analize_data("MLA409431")
	}
}
