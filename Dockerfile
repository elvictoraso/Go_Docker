# Cambia la versión vieja (1.21) por la de tu go.mod (1.26)
FROM golang:1.26-alpine AS builder

WORKDIR /app

# Copiar archivos de dependencias
COPY go.mod ./
# Si ya tienes go.sum, descomenta la siguiente línea:
# COPY go.sum ./

RUN go mod download

# Copiar el resto del código fuente
COPY . .

# Compilar de forma estática
RUN CGO_ENABLED=0 GOOS=linux go build -o password-checker .

# Etapa final optimizada
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/password-checker .
COPY --from=builder /app/templates ./templates
# (Asegúrate de copiar carpetas de assets/estáticos si los usas)

EXPOSE 8080
CMD ["./password-checker"]
