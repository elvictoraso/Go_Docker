# Etapa 1: Construcción
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copiamos todo el código fuente
COPY . .

# Inicializamos el módulo de Go de forma segura y compilamos
RUN go mod init password-checker || true && \
    go mod tidy && \
    CGO_ENABLED=0 GOOS=linux go build -o password-checker .

# Etapa 2: Imagen final ligera
FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copiamos el binario de la etapa anterior
COPY --from=builder /app/password-checker .

EXPOSE 8080

CMD ["./password-checker"]
