// Autor: Javier Gonzalez
// Fecha 20-11-2017
// email: mobile.javierg@gmail.como

package clientApi

import (
	"sort"
)

type ResponseAPI struct {
	Max       float64 `json:"max"`
	Seggested float64 `json:"suggested"`
	Min       float64 `json:"min"`
}

type Categories struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type Paging struct {
	Total           int `json:"total"`
	Offset          int `json:"offset"`
	Limit           int `json:"limit"`
	Primary_results int `"primary_results"`
}

type Result struct {
	Id            int     `json:"id"`
	Price         float64 `json:"price"`
	Sold_quantity int     `json:"sold_quantity"`
}

type Search struct {
	first      bool
	Paginado   Paging   `json:"paging"`
	Resultados []Result `json:"results"`
}

// implemento estas funciones para poder ordenar un array struct
func (s Search) Len() int {
	return len(s.Resultados)
}
func (s Search) Swap(i, j int) {
	s.Resultados[i], s.Resultados[j] = s.Resultados[j], s.Resultados[i] //Aqui las dudas
}
func (s Search) Less(i, j int) bool {
	return s.Resultados[i].Price < s.Resultados[j].Price
}

func GetEstadistics(result Search) (max float64, min float64, mediada float64) {

	//ordenamos el array
	sort.Sort(result)

	var _Max float64 = 0
	var _Min float64 = 0
	var _Media float64 = 0

	var _acum float64

	//////////////////////////////////////////
	//recorro y analizo el resultado obtenido

	_Min = result.Resultados[0].Price
	_Max = result.Resultados[0].Price

	for i := 0; i <= len(result.Resultados)-1; i++ {

		// tomo las muestras como valida, si vendieron aunque sea un articulo
		// ya que hay publicaciones con valores irrales, por ejemplo: 111111111.11, o 123456789.00
		if result.Resultados[i].Sold_quantity > 0 {

			_acum += result.Resultados[i].Price

			//guardo el menor
			if _Min > result.Resultados[i].Price {
				_Min = result.Resultados[i].Price
			}

			//guardo el mayor
			if _Max < result.Resultados[i].Price {
				_Max = result.Resultados[i].Price
			}

		}

	}

	//obtengo la media aritmetica
	_Media = _acum / float64(len(result.Resultados))

	return _Max, _Min, _Media

}
