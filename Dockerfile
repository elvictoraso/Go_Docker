# Etapa 1: Construcción
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o password-checker .

# Etapa 2: Imagen final optimizada
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/password-checker .
EXPOSE 8080
CMD ["./password-checker"]
