# Estágio de build
FROM golang:1.23-alpine AS builder
RUN apk --no-cache add tzdata gcc musl-dev sqlite-dev
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . .
# RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags='-w -s -extldflags "-static"' -o ibge-service ./cmd/ibge-api
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -a -ldflags='-w -s -linkmode external -extldflags "-static"' -o ibge-service ./cmd/ibge-api

# Estágio final de produção
FROM alpine:latest
RUN apk --no-cache add tzdata wget
WORKDIR /app
COPY --from=builder /app/ibge-service .
COPY --from=builder /app/data/ibge.db /app/data/
COPY --from=builder /app/docs/ /app/docs/
# COPY ./data/ibge.db /app/data/ibge.db

# Criar usuário não-root (melhor segurança)
RUN adduser -D -s /bin/sh appuser && \
    chown -R appuser:appuser /app
USER appuser
# Configurações de ambiente
ENV GO_ENV=production
ENV LOG_LEVEL=warn
# Configurações de rate limiting
ENV RATE_LIMIT=50
ENV RATE_LIMIT_WINDOW=1m
# Expor a porta que a aplicação usará (deve corresponder à SERVER_PORT)
EXPOSE 9080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:9080/health || exit 1

# Comando para rodar a aplicação
ENTRYPOINT ["./ibge-service"]