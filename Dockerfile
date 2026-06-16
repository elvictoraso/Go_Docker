# Etapa 1: Construcción
FROM golang:1.21-alpine AS builder
WORKDIR /app

# Copiamos directamente el código fuente y el módulo
COPY . .

# Compilamos la aplicación de forma estática
RUN CGO_ENABLED=0 GOOS=linux go build -o password-checker .

# Etapa 2: Imagen final optimizada
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/

# Traemos el binario compilado desde la etapa anterior
COPY --from=builder /app/password-checker .

EXPOSE 8080
CMD ["./password-checker"]
