# Mercadolibre - Golang

enunciado [aqui](https://github.com/mobilejavierg/mercadolibre/blob/master/enunciado.md)

Se debe analizar el precio de los artículos para una categoría, como es muy costoso analizar todo, utilice el método de Muestreo Aleatorio Simple para determinar el tamaño de la cantidad de datos a procesar, y así poder calcular la media aritmética. Al momento de analizar el precio de los artículos tuve que tener en cuenta el atributo Sold_quantity>0 (al menos un artículo vendido) ya que había muchas publicaciones con valores irreales.

Este ejercicio sirvió para practicar los temas vistos en el curso y tomar como desafío aprender un nuevo lenguaje.

Para poder analizar más rápido los artículos tuve que utilizar goroutines, aprovechando las bondades de la “concurrencia” y/o “paralelismo” que ofrece este lenguaje moderno, así como también el uso de channels para acumular los resultados y wait group que me sirvió para la sincronización.

A parte de Golang otro mundo nuevo fue implementar la solución en Google Cloud, ya que tuve refactorizar el código para adaptar la solución al appengine.

La primer situación surgió en agregar el módulo de "google.golang.org/appengine" y modificar mi func main(), agregando init().

Utilice gin-gonic para rutear los get’s, donde surgió las 2da situación había conflicto con el Listening Port ya que estaba siendo utilizada por gapp. Se solucionó vinculando gin con net/http y acoplándome al handle de appengine “http.Handle("/", r)”

La siguiente situación fue el conflicto generado por appengine con los Get “externos” hacia las api’s de mercadolibre, tuve que utilizar el modulo "google.golang.org/appengine/urlfetch" y hacer burbujear el http.request desde gin-gonic hacia las funciones asincrónicas donde consumo las apis.

Consumir el API de obtener artículos por categoría es muy costoso en tiempos, tuve modificar el context generado alargando el time out:

ctx := appengine.NewContext(req)

ctx, cancel := context.WithTimeout(ctx, 20*time.Second)

Al aplicar la formula de Muestro Aleatorio Simple, tuve que subir la tasa de error al 5%, para asi poder dismunir la cantidad de articulos a analizar, para bajar la latencia en la respuesta.



##Mejoras
- Guardaria los datos en una base de datos, para poder analizar los mismo de manera offline y mejorar la velocidad de respuesta, ya que actualmente analiza los datos en forma "on the fly".
-
