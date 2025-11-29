# Generador y Deploy de Funciones Lambda para AWS

Este directorio contiene herramientas para generar y desplegar funciones Lambda automáticamente basándose en `functions.json`.

## Estructura

```
functions/
├── main.go              # CLI principal que controla generate y deploy
├── go.mod               # Módulo Go para el CLI
├── deploy.sh            # Script helper para deploy
├── utils/               # Utilidades compartidas
│   ├── generate.go     # Lógica de generación
│   ├── deploy.go       # Lógica de deploy
│   └── types.go        # Estructuras compartidas
├── templates/          # Templates para generar código
│   ├── main.go.tmpl   # Template para main.go
│   └── go.mod.tmpl    # Template para go.mod
└── [módulo]/           # Funciones generadas
    └── [función]/
        ├── main.go
        ├── go.mod
        └── bin/        # Directorio para artefactos compilados
```

## Uso del CLI

### Generar Funciones

Para generar todas las funciones desde `functions.json`:

```bash
cd src/infrastructure/clouds/aws/functions
go run main.go generate
```

O usando alias:
```bash
go run main.go gen
go run main.go g
```

### Desplegar Funciones

**Desplegar una función específica:**
```bash
# Usando el CLI
go run main.go deploy health-check

# O usando el script helper
./deploy.sh health-check
```

**Desplegar todas las funciones:**
```bash
# Usando el CLI
go run main.go deploy

# O usando el script helper
./deploy.sh
```

### Ayuda

```bash
go run main.go help
go run main.go --help
go run main.go -h
```

## Proceso de Deploy

El comando `deploy` automatiza los siguientes pasos para cada función:

1. **go mod tidy** - Limpia y actualiza dependencias
2. **Compilación** - Compila con:
   ```bash
   GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o bin/bootstrap -tags lambda.norpc main.go
   ```
3. **Empaquetado** - Crea un zip con el binario:
   ```bash
   zip -j bin/[function-name].zip bin/bootstrap
   ```
4. **Deploy** - Actualiza la función Lambda:
   ```bash
   aws lambda update-function-code \
     --function-name [function-name] \
     --zip-file fileb://bin/[function-name].zip \
     --region [region]
   ```

## Configuración

El script usa variables de entorno con valores por defecto:

- `PROJECT_NAME` (default: "go-project-skeleton")
- `ENVIRONMENT` (default: "development")
- `AWS_REGION` (default: "us-east-1")

Puedes configurarlas antes de ejecutar:

```bash
export PROJECT_NAME="mi-proyecto"
export ENVIRONMENT="production"
export AWS_REGION="us-west-2"
./deploy.sh health-check
```

## Nombres de Funciones Lambda

El script calcula el nombre de la función Lambda usando la misma lógica que Terraform:

- `name_prefix = "${project_name_clean}-${env_prefix}"`
- `function_name = "${name_prefix}-${function_name_clean}"`

Ejemplo:
- Proyecto: "go-project-skeleton"
- Ambiente: "development"
- Función: "health-check"
- Resultado: "goprojectskeleton-dev-healthcheck"

## Requisitos

- Go 1.25+
- AWS CLI configurado con credenciales válidas
- Permisos para actualizar funciones Lambda
- `zip` instalado en el sistema

## Notas

- El script preserva el directorio `bin/` si ya existe
- Los archivos `main.go` y `go.mod` se regeneran completamente al ejecutar `generate.go`
- El deploy solo actualiza el código, no crea nuevas funciones (deben existir en Lambda)
- Asegúrate de tener `functions.json` en el directorio padre (`../functions.json`)
