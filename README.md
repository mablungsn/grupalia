# API

Este repositorio contiene la API de la aplicación web.

Toda la aplicación tiene los tenedores necesarios para funcionar y tiene un seed que permite poblar la base de datos con algunos valores

Notar que se requiere docker y docker compose instalado

Hacer build
```sh
docker compose build
```

Correr 
```sh
docker compose up
```

Borrar (sirve para volver a correr build)
```sh
docker compose down
```

Otros comando utiles de la API

```sh
curl -X POST --url http://localhost:8080/authentication/login -d '{"email": "prueba1@hola.cl", "password":"pass1"}'
curl -X POST --url http://localhost:8080/authentication/login -d '{"email": "prueba2@hola.cl", "password":"pass2"}'
curl --request GET --url http://localhost:8080/v1 --header 'Authorization: Bearer <acces token de request anterior>' -v
```





