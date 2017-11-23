# Enunciado
# Ejercicio Integrador - Sugeridor de precios

En Meli todo los días miles de usuarios realizan millones de publicaciones y uno de los principales desafíos a los que se enfrentan al crear su publicación es entender cuál sería el precio ideal para su producto. Un precio alto sería poco competitivo, un precio bajo reduce sus ganancias. Tu desafío como nuevo Golang developer en nuestro team es facilitarle la toma de esta decisión creando una API que le sugiera el precio ideal para su producto.

Detalles:

Crear una api REST en Golang ;) que retorne para una determinada categoría de productos un JSON como respuesta con el siguiente formato:
{ 
"max":10, 
"suggested":5, 
"min":1 
}

sobre el siguiente recurso:

/categories/$ID/prices

ejemplo:

curl -X GET “http://mydomain.com/categories/MLA3530/prices”

Realizar test (white/black box test) y benchmark del servicio, recuerden probar concurrencia! Una cobertura mayor al 80% del código es aceptable ;)
Hostear el código a un repo git y la API en un Cloud Computing público.
