# web-simplify
![technology Go](https://img.shields.io/badge/technology-go-blue.svg)
![TCP/IP Communication](https://img.shields.io/badge/communication-TCP%2FIP-orange.svg)
![HTTP](https://img.shields.io/badge/protocol-HTTP-lightgrey.svg)
![Viper](https://img.shields.io/badge/configuration-viper-green.svg)

## Descripción

web-simplify es una librería de Go que simplifica la creación de aplicaciones web, integrando herramientas como conexión a bases de datos, router con Chi, cliente REST, manejo de errores, Viper para configuraciones, y documentación con Docsify y Swagger. Los servicios y componentes se crean automáticamente desde el archivo de configuración.

## 🚀 Instalación
Para instalar Web-Simplify, ejecuta:
```sh
 go get github.com/skolldire/web-simplify
```

## 📌 Características

### Servers:

1. **Servidor TCP/IP**
    * **Package:** GoTCPServer
    * **Descripción:** Proporciona una herramienta para crear servidores TCP/IP, permitiendo la creación de sockets de comunicación entre sistemas y facilitando la transmisión de datos en tiempo real a través de la red.

### Clients:

1. **Conexión a bases de datos**
    * **Package:** db_connect
    * **Descripción:** Facilita la integración y gestión de bases de datos utilizando XORM como capa ORM, permitiendo acceder y manipular datos de manera eficiente.
2. **Router con Chi**
    * **Package:** simple_router
    * **Descripción:** Gestiona las rutas de las aplicaciones RESTful utilizando Chi, un router ligero y eficiente para Go.

### Utilities:

1. **Manejo de errores**
    * **Package:** error_handler
    * **Descripción:** Incluye un manejador de errores (error wrapper) que permite un manejo limpio y eficiente de excepciones en la aplicación.
2. **Utilidades para manejo de archivos**
    * **Package:** file_utils
    * **Descripción:** Herramientas para la gestión eficiente de archivos dentro de la aplicación, como lectura, escritura y manipulación de archivos.
3. **Ejecución de tareas en paralelo**
    * **Package:** task_runner
    * **Descripción:** Facilita la creación y gestión de goroutines, permitiendo la ejecución de tareas de forma concurrente y eficiente.
4. **Gestión de perfiles de aplicación**
    * **Package:** app_profile
    * **Descripción:** Permite gestionar perfiles de configuración para diferentes entornos (desarrollo, producción, etc.), facilitando la segregación de ambientes y la creación de instancias.
5. **Conversión de datos**
    * **Package:** data_converter
    * **Descripción:** Proporciona funcionalidades para la conversión rápida entre diferentes tipos de datos y estructuras.
6. **Manejo de logs personalizable**
    * **Package:** simple_logger
    * **Descripción:** Implementa Logrus para el manejo de logs, ofreciendo personalización avanzada de los logs generados por la aplicación, con distintos niveles de severidad y formatos.

## 🛠️ Uso Básico
```go
package main

import (
    "github.com/skolldire/web-simplify/pkg/simplify"
)

func main() {
    app := simplify.New()
    app.Start()
}
```

## ✅ Requisitos
- Go 1.20+
- Dependencias como `chi`, `logrus`, `xorm` y `viper` instaladas

## 🔍 Documentación
Puedes encontrar la documentación completa en [pkg.go.dev](https://pkg.go.dev/github.com/skolldire/web-simplify)

También puedes generar documentación localmente ejecutando:
```sh
godoc -http=:6060
```
Y accediendo a `http://localhost:6060/pkg/github.com/skolldire/web-simplify/`

## 🧪 Pruebas
Para ejecutar las pruebas del proyecto, usa:
```sh
go test ./pkg/... -cover
```
Esto mostrará la cobertura de pruebas y validará el correcto funcionamiento de la librería.

## 🧑‍💻 Contribuciones
¡Las contribuciones son bienvenidas! Consulta [CONTRIBUTING.md](./CONTRIBUTING.md) para más detalles.

Para reportar errores o sugerencias, abre un Issue en [GitHub](https://github.com/skolldire/web-simplify/issues).

## 📦 Publicación
Si deseas usar una versión específica, puedes instalarla con:
```sh
go get github.com/skolldire/web-simplify@v1.0.0
```

## 📜 Licencia
Este proyecto está bajo la licencia MIT. Consulta el archivo [LICENSE](./LICENSE) para más información.
