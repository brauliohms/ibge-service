#!/bin/bash

# Script para criar o binário final do projeto

set -e

# Opcional: crie um diretório para os binários
mkdir -p bin

# Compile o projeto. O -o especifica o nome e o local do arquivo de saída.
go build -o ./bin/ibge-service ./cmd/ibge-api/main.go
