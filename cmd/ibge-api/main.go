package main

import (
	"fmt"
	"log"
	"net/http"
	"runtime"

	"github.com/brauliohms/ibge-service/docs"
	httphandler "github.com/brauliohms/ibge-service/internal/adapter/http"
	"github.com/brauliohms/ibge-service/internal/adapter/repository/memory"
	"github.com/brauliohms/ibge-service/internal/adapter/repository/sqlite"
	"github.com/brauliohms/ibge-service/internal/usecase"
	"github.com/brauliohms/ibge-service/pkg/config"
)

// @title           API de Dados do IBGE
// @version         1.0
// @description     Este é um microserviço para consulta de estados e cidades do Brasil, baseado nos dados do IBGE.

// @contact.name   API IBGE Support
// @contact.url    http://www.exemplo.com/support
// @contact.email  contato@integradocs.com.br

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @BasePath  /api/v1

func main() {
	// Otimizar para alta concorrência
	runtime.GOMAXPROCS(runtime.NumCPU())
	// 1. Carregar configurações
	cfg := config.Load()
	log.Println("Configurações carregadas")

	// 2. Inicializar o repositório fonte (PostgreSQL)
	// Este repositório será usado APENAS para a carga inicial.
	// ibgeRepo, err := postgres.NewPostgresRepository(cfg.PostgresDSN)
	// if err != nil {
	// 	log.Fatalf("Falha ao conectar com o PostgreSQL: %v", err)
	// }
	// log.Println("Conexão com PostgreSQL estabelecida")

	ibgeRepo, err := sqlite.NewSQLiteRepository(cfg.SqliteDSN)
	if err != nil {
		log.Fatalf("Falha ao conectar com o SQLite: %v", err)
	}

	// 3. Inicializar o repositório em memória, usando o PostgreSQL como fonte.
	// Esta é a mágica do cache no startup!
	memoryRepo, err := memory.NewMemoryRepository(ibgeRepo)
	if err != nil {
		log.Fatalf("Falha ao carregar dados para o repositório em memória: %v", err)
	}
	log.Println("Repositório em memória populado com sucesso")
	// Neste ponto, a conexão com o Postgres poderia até ser fechada se não fosse mais necessária.

	// 4. Injetar o repositório EM MEMÓRIA no caso de uso.
	// A aplicação agora só falará com o cache, sendo extremamente rápida.
	ibgeUseCase := usecase.NewIBGEUseCase(memoryRepo)

	// 5. Injetar o caso de uso nos handlers HTTP.
	ibgeHandler := httphandler.NewIBGEHandler(ibgeUseCase)

	// 6. Configurar o roteador.
	router := httphandler.SetupRouter(ibgeHandler)

	docs.SwaggerInfo.Title = "API de Dados do IBGE"
	docs.SwaggerInfo.Description = "Microserviço para consulta de estados e cidades do Brasil."
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = fmt.Sprintf("localhost:%s", cfg.ServerPort)
	docs.SwaggerInfo.BasePath = "/api/v1"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

	// 7. Iniciar o servidor HTTP.
	serverAddr := fmt.Sprintf(":%s", cfg.ServerPort)
	log.Printf("Servidor iniciado e ouvindo em http://localhost%s", serverAddr)
	if err := http.ListenAndServe(serverAddr, router); err != nil {
		log.Fatalf("Falha ao iniciar o servidor HTTP: %v", err)
	}
}
