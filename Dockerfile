# Imagen base oficial de Go
FROM golang:1.24.5

# Crear directorio de trabajo
WORKDIR /app

# Copiar solo los archivos de dependencias y descargarlas
COPY go.mod go.sum ./
RUN go mod download

# Copiar el resto del código
COPY . .

# Compilar el binario y guardarlo en /app/main
RUN go build -o /main ./cmd/server

# Mostrar los archivos en / y /app (debug)
RUN ls -lah / && ls -lah /app

# Ejecutar el binario usando la ruta absoluta
CMD ["/main"]
