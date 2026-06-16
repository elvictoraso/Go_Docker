# Etapa 1: Construcción
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copiamos primero los archivos del módulo para aprovechar la caché de Docker
COPY go.mod ./
# Si el comando go mod tidy te generó un archivo go.sum, descomenta la siguiente línea:
# COPY go.sum ./

RUN go mod download

# Copiamos el resto del código fuente
COPY . .

# Compilamos la aplicación de forma estática
RUN CGO_ENABLED=0 GOOS=linux go build -o password-checker .

# Etapa 2: Imagen final ligera
FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copiamos el binario compilado de la etapa anterior
COPY --from=builder /app/password-checker .

EXPOSE 8080

CMD ["./password-checker"]
