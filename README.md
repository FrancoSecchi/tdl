# Teoría del lenguaje

## Inicialización del proyecto

Para poder levantar correctamente el proyecto, es necesario que tengan instalado el lenguaje Go en su máquina. [Página para descargar](https://go.dev/dl/)

Luego de tenerlo descargado. Hay que ubicarse en la carpeta raíz del proyecto y ejecutar el siguiente comando en consola:

```
go run main.go
```

Si hace cambios en los archivos .go, tiene que ejecutar primero:

```
go build main.go
```

Y despues:

```
go run main.go
```

## ¿De qué se trata?

Este proyecto es un chat implementado en Go, Html, CSS y Js. En el cual varias personas se pueden registrar y loguear y enviarse mensajes entre ellos y entre un chat grupal en donde se encuentran todos.

Para poder implementarlo se tuvo que investigar sobre el funcionamiento de Go, las goroutines, y sobre todo los websockets.


