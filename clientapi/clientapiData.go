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

	_Min = result.Resultados[0].Price
	_Max = result.Resultados[len(result.Resultados)-1].Price

	for i := 0; i <= len(result.Resultados)-1; i++ {
		_acum += result.Resultados[i].Price
	}

	_Media = _acum / float64(len(result.Resultados))

	return _Max, _Min, _Media

}
