package main

import (
	"database/sql"
	"flag"
	"log"

	"github.com/brauliohms/ibge-service/internal/seed"

	// Importar drivers de banco de dados
	_ "github.com/lib/pq"           // PostgreSQL
	_ "github.com/mattn/go-sqlite3" // SQLite
	// _ "github.com/go-sql-driver/mysql" // MySQL (se necessário)
)

func main() {
	var (
		driverName = flag.String("driver", "postgres", "Nome do driver do banco de dados (postgres, sqlite3, mysql)")
		dsn        = flag.String("dsn", "", "Data Source Name para conexão com o banco")
		dataDir    = flag.String("data-dir", "data", "Diretório contendo os arquivos JSON")
	)
	flag.Parse()

	if *dsn == "" {
		log.Fatal("DSN é obrigatório. Use -dsn para especificar a string de conexão")
	}

	// Normalizar nome do driver
	driver := *driverName
	if driver == "sqlite" {
		driver = "sqlite3"
	}

	log.Printf("Conectando ao banco de dados com driver: %s", driver)

	// Conectar ao banco de dados
	db, err := sql.Open(driver, *dsn)
	if err != nil {
		log.Fatalf("Erro ao conectar ao banco de dados: %v", err)
	}
	defer db.Close()

	// Testar conexão
	if err := db.Ping(); err != nil {
		log.Fatalf("Erro ao testar conexão com o banco: %v", err)
	}

	// Executar seed
	seeder := seed.NewSeeder(db, driver)
	if err := seeder.Run(*dataDir); err != nil {
		log.Fatalf("Erro durante o seed: %v", err)
	}

	log.Println("Seed executado com sucesso!")
}

// Compilar e executar o seed
// go run cmd/seed/main.go -driver="$DRIVER" -dsn="$DB_DSN" -data-dir="$DATA_DIR"
