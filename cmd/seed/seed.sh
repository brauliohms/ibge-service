#!/bin/bash

# Script para executar o seed do banco de dados

set -e

# Configurações padrão
DRIVER=${DB_DRIVER:-"sqlite3"}
DATA_DIR=${DATA_DIR:-"data"}
DB_DSN=${DB_DSN:-"./data/ibge.db"}

# Função para mostrar ajuda
show_help() {
  echo "Uso: $0 [postgres|sqlite|mysql]"
  echo ""
  echo "Variáveis de ambiente:"
  echo "  DB_DRIVER - Driver do banco (postgres, sqlite3, mysql)"
  echo "  DB_DSN    - String de conexão do banco"
  echo "  DATA_DIR  - Diretório dos arquivos JSON (padrão: data)"
  echo ""
  echo "Exemplos:"
  echo "  # PostgreSQL"
  echo "  DB_DSN='postgres://user:pass@localhost/dbname?sslmode=disable' $0 postgres"
  echo ""
  echo "  # SQLite"
  echo "  DB_DSN='./data/ibge.db' $0 sqlite"
  echo ""
  echo "  # MySQL"
  echo "  DB_DSN='user:pass@tcp(localhost:3306)/dbname' $0 mysql"
}

# Verificar argumentos
if [ $# -gt 0 ]; then
  DRIVER=$1
fi

if [ "$DRIVER" = "help" ] || [ "$DRIVER" = "-h" ] || [ "$DRIVER" = "--help" ]; then
  show_help
  exit 0
fi

# Verificar se DSN foi fornecido
if [ -z "$DB_DSN" ]; then
  echo "Erro: DB_DSN não foi definido"
  echo ""
  show_help
  exit 1
fi

# Verificar se os arquivos JSON existem
if [ ! -f "$DATA_DIR/cidades-ibge-uf.json" ]; then
  echo "Erro: Arquivo $DATA_DIR/cidades-ibge-uf.json não encontrado"
  exit 1
fi

if [ ! -f "$DATA_DIR/municipios-TOM-IBGE.json" ]; then
  echo "Erro: Arquivo $DATA_DIR/municipios-TOM-IBGE.json não encontrado"
  exit 1
fi

echo "Executando seed com:"
echo "  Driver: $DRIVER"
echo "  DSN: $DB_DSN"
echo "  Data Dir: $DATA_DIR"
echo ""

# Compilar e executar o seed
go run cmd/seed/main.go -driver="$DRIVER" -dsn="$DB_DSN" -data-dir="$DATA_DIR"
