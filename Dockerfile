# Etapa 1: Compilación
FROM golang:1.26-alpine AS builder

WORKDIR /app

# Copiar archivos de dependencias
COPY go.mod ./
RUN go mod download

# Copiar el resto del código fuente
COPY . .

# Compilar de forma estática
RUN CGO_ENABLED=0 GOOS=linux go build -o password-checker .

# Etapa 2: Imagen final optimizada
FROM alpine:latest
WORKDIR /app

# Copiamos el ejecutable desde el builder
COPY --from=builder /app/password-checker .

# CORRECCIÓN: Eliminamos o comentamos la línea de templates porque no existe en tu proyecto
# COPY --from=builder /app/templates ./templates

EXPOSE 8080
CMD ["./password-checker"]
