# Go Generator Phrases - LocalServer

Aplicacion web local en Go que genera frases aleatorias y las muestra como imagenes PNG en una galeria servida por HTTP.

La rama `LocalServer` conserva su objetivo original: ejecutar un servidor local con interfaz web. A diferencia de `main`, aqui el foco no es la libreria publica sino la experiencia web local.

## Caracteristicas

- Servidor HTTP local en `http://localhost:8080`.
- Generacion dinamica de imagenes PNG con frases aleatorias.
- Carga incremental de imagenes desde la interfaz web.
- Diccionarios embebidos con `embed` para evitar rutas fragiles.
- Ajuste automatico de fuente para mejorar legibilidad de cada imagen.
- Validaciones basicas del endpoint y manejo de errores mas claro.

## Requisitos

- Go 1.22 o superior.

## Ejecutar

```bash
go run .
```

Luego abre:

```text
http://localhost:8080
```

## Endpoints

| Ruta | Metodo | Descripcion |
| --- | --- | --- |
| `/` | `GET` | Sirve la interfaz HTML |
| `/load-more-images` | `GET` | Devuelve una lista JSON de imagenes generadas |
| `/img/*` | `GET` | Sirve los PNG generados |

El endpoint `/load-more-images` acepta `count` como query param. Rango permitido: `1` a `50`.

Ejemplo:

```text
/load-more-images?count=12
```

## Arquitectura

```mermaid
flowchart TD
    A[Browser] --> B[HTTP Server main.go]
    B --> C[templates/index.html]
    B --> D[/load-more-images]
    D --> E[image.GenerateImages]
    E --> F[generator.New]
    F --> G[Diccionarios embebidos]
    E --> H[Render PNG con gg]
    H --> I[img/*.png]
    I --> A
```

## Estructura

```text
.
|-- fonts/
|-- generator/
|   |-- dictionaries/
|   |-- generator.go
|   `-- generator_test.go
|-- image/
|   `-- image.go
|-- templates/
|   `-- index.html
|-- main.go
|-- go.mod
`-- README.md
```

## Desarrollo

Ejecutar pruebas:

```bash
go test ./...
```

## Mejoras integradas desde `Images`

- Formateo de texto mas seguro sin `strings.Title`.
- Renderizado de imagen con mejor ajuste de fuente.
- Creacion automatica del directorio de salida.
- Generador desacoplado de rutas fisicas usando `embed`.
- Menos duplicacion y mejor manejo de errores en servidor y generacion.
