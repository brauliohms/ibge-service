package seed

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

// CidadeIBGE representa a estrutura do arquivo cidades-ibge-uf.json
type CidadeIBGE struct {
	UF             string `json:"uf"`
	UFNome         string `json:"ufNome"`
	CodigoUF       int    `json:"codigoUf"`
	Codigo         int    `json:"codigo"`
	Nome           string `json:"nome"`
	Microregiao    int    `json:"microrregiao"`
	RegiaoImediata int    `json:"regiaoImediata"`
}

// MunicipioTOM representa a estrutura do arquivo municipios-TOM-IBGE.json
type MunicipioTOM struct {
	CodigoMunicipioIBGE int `json:"CODIGO-MUNICIPIO-IBGE"`
	CodigoTOM           int `json:"CODIGO-MUNICIPIO-TOM"`
}

// Seeder gerencia o processo de seed do banco de dados
type Seeder struct {
	db         *sql.DB
	driverName string
}

// NewSeeder cria uma nova instância do seeder
func NewSeeder(db *sql.DB, driverName string) *Seeder {
	return &Seeder{
		db:         db,
		driverName: driverName,
	}
}

// Run executa o processo completo de seed
func (s *Seeder) Run(dataDir string) error {
	log.Println("Iniciando processo de seed...")

	// 1. Criar tabelas
	if err := s.createTables(); err != nil {
		return fmt.Errorf("erro ao criar tabelas: %w", err)
	}

	// 2. Carregar dados dos arquivos JSON
	cidades, err := s.loadCidadesData(filepath.Join(dataDir, "cidades-ibge-uf.json"))
	if err != nil {
		return fmt.Errorf("erro ao carregar dados das cidades: %w", err)
	}

	tomData, err := s.loadTOMData(filepath.Join(dataDir, "municipios-TOM-IBGE.json"))
	if err != nil {
		return fmt.Errorf("erro ao carregar dados TOM: %w", err)
	}

	// 3. Popular estados
	if err := s.seedEstados(cidades); err != nil {
		return fmt.Errorf("erro ao popular estados: %w", err)
	}

	// 4. Popular cidades
	if err := s.seedCidades(cidades, tomData); err != nil {
		return fmt.Errorf("erro ao popular cidades: %w", err)
	}

	log.Println("Processo de seed concluído com sucesso!")
	return nil
}

// createTables cria as tabelas no banco de dados
func (s *Seeder) createTables() error {
	log.Println("Criando tabelas...")

	var createEstadosSQL, createCidadesSQL, createIndexSQL string

	switch s.driverName {
	case "sqlite3":
		createEstadosSQL = `
		CREATE TABLE IF NOT EXISTS estados (
			codigo_ibge INTEGER PRIMARY KEY,
			nome VARCHAR(30) NOT NULL UNIQUE,
			sigla CHAR(2) NOT NULL UNIQUE
		);`

		createCidadesSQL = `
		CREATE TABLE IF NOT EXISTS cidades (
			codigo_ibge INTEGER PRIMARY KEY,
			nome VARCHAR(50) NOT NULL,
			codigo_tom INT,
			micro_regiao INT,
			regiao_imediata INT,
			estado_codigo_ibge INTEGER NOT NULL,
			FOREIGN KEY(estado_codigo_ibge) REFERENCES estados(codigo_ibge)
		);`

		createIndexSQL = `
		CREATE INDEX IF NOT EXISTS idx_cidades_por_estado ON cidades(estado_codigo_ibge);`

	case "postgres":
		createEstadosSQL = `
		CREATE TABLE IF NOT EXISTS estados (
			codigo_ibge INT PRIMARY KEY,
			nome VARCHAR(30) NOT NULL UNIQUE,
			sigla CHAR(2) NOT NULL UNIQUE
		);`

		createCidadesSQL = `
		CREATE TABLE IF NOT EXISTS cidades (
			codigo_ibge INT PRIMARY KEY,
			nome VARCHAR(50) NOT NULL,
			codigo_tom INT,
			micro_regiao INT,
			regiao_imediata INT,
			estado_codigo_ibge INT NOT NULL,
			CONSTRAINT fk_estado
				FOREIGN KEY(estado_codigo_ibge) 
				REFERENCES estados(codigo_ibge)
		);`

		createIndexSQL = `
		CREATE INDEX IF NOT EXISTS idx_cidades_por_estado ON cidades(estado_codigo_ibge);`

	default: // MySQL e outros
		createEstadosSQL = `
		CREATE TABLE IF NOT EXISTS estados (
			codigo_ibge INT PRIMARY KEY,
			nome VARCHAR(30) NOT NULL UNIQUE,
			sigla CHAR(2) NOT NULL UNIQUE
		);`

		createCidadesSQL = `
		CREATE TABLE IF NOT EXISTS cidades (
			codigo_ibge INT PRIMARY KEY,
			nome VARCHAR(50) NOT NULL,
			codigo_tom INT,
			micro_regiao INT,
			regiao_imediata INT,
			estado_codigo_ibge INT NOT NULL,
			CONSTRAINT fk_estado
				FOREIGN KEY(estado_codigo_ibge) 
				REFERENCES estados(codigo_ibge)
		);`

		createIndexSQL = `
		CREATE INDEX idx_cidades_por_estado ON cidades(estado_codigo_ibge);`
	}

	// Executar SQLs
	if _, err := s.db.Exec(createEstadosSQL); err != nil {
		return fmt.Errorf("erro ao criar tabela estados: %w", err)
	}

	if _, err := s.db.Exec(createCidadesSQL); err != nil {
		return fmt.Errorf("erro ao criar tabela cidades: %w", err)
	}

	if _, err := s.db.Exec(createIndexSQL); err != nil {
		log.Printf("Aviso ao criar índice (pode já existir): %v", err)
	}

	log.Println("Tabelas criadas com sucesso!")
	return nil
}

// loadCidadesData carrega os dados do arquivo cidades-ibge-uf.json
func (s *Seeder) loadCidadesData(filePath string) ([]CidadeIBGE, error) {
	log.Printf("Carregando dados de cidades de: %s", filePath)

	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("erro ao ler arquivo: %w", err)
	}

	var cidades []CidadeIBGE
	if err := json.Unmarshal(data, &cidades); err != nil {
		return nil, fmt.Errorf("erro ao fazer unmarshal do JSON: %w", err)
	}

	log.Printf("Carregadas %d cidades", len(cidades))
	return cidades, nil
}

// loadTOMData carrega os dados do arquivo municipios-TOM-IBGE.json
func (s *Seeder) loadTOMData(filePath string) (map[int]string, error) {
	log.Printf("Carregando dados TOM de: %s", filePath)

	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("erro ao ler arquivo: %w", err)
	}

	var tomList []MunicipioTOM
	if err := json.Unmarshal(data, &tomList); err != nil {
		return nil, fmt.Errorf("erro ao fazer unmarshal do JSON: %w", err)
	}

	// Converter para map para busca rápida
	tomMap := make(map[int]string)
	for _, tom := range tomList {
		tomMap[tom.CodigoMunicipioIBGE] = fmt.Sprintf("%d", tom.CodigoTOM)
	}

	log.Printf("Carregados %d códigos TOM", len(tomMap))
	return tomMap, nil
}

// seedEstados popula a tabela de estados
func (s *Seeder) seedEstados(cidades []CidadeIBGE) error {
	log.Println("Populando estados...")

	// Extrair estados únicos
	estadosMap := make(map[int]CidadeIBGE)
	for _, cidade := range cidades {
		if _, exists := estadosMap[cidade.CodigoUF]; !exists {
			estadosMap[cidade.CodigoUF] = cidade
		}
	}

	// Preparar statement baseado no driver
	var stmt *sql.Stmt
	var err error

	switch s.driverName {
	case "postgres":
		stmt, err = s.db.Prepare(`
			INSERT INTO estados (codigo_ibge, nome, sigla) 
			VALUES ($1, $2, $3) 
			ON CONFLICT (codigo_ibge) DO NOTHING
		`)
	case "sqlite3":
		stmt, err = s.db.Prepare(`
			INSERT OR IGNORE INTO estados (codigo_ibge, nome, sigla) 
			VALUES (?, ?, ?)
		`)
	default:
		// MySQL e outros
		stmt, err = s.db.Prepare(`
			INSERT IGNORE INTO estados (codigo_ibge, nome, sigla) 
			VALUES (?, ?, ?)
		`)
	}

	if err != nil {
		return fmt.Errorf("erro ao preparar statement: %w", err)
	}
	defer stmt.Close()

	// Inserir estados
	count := 0
	for _, estado := range estadosMap {
		if _, err := stmt.Exec(estado.CodigoUF, estado.UFNome, estado.UF); err != nil {
			return fmt.Errorf("erro ao inserir estado %s: %w", estado.UF, err)
		}
		count++
	}

	log.Printf("Processados %d estados", count)
	return nil
}

// seedCidades popula a tabela de cidades
func (s *Seeder) seedCidades(cidades []CidadeIBGE, tomData map[int]string) error {
	log.Println("Populando cidades...")

	// Preparar statement baseado no driver
	var stmt *sql.Stmt
	var err error

	switch s.driverName {
	case "postgres":
		stmt, err = s.db.Prepare(`
			INSERT INTO cidades (codigo_ibge, nome, codigo_tom, micro_regiao, regiao_imediata, estado_codigo_ibge) 
			VALUES ($1, $2, $3, $4, $5, $6) 
			ON CONFLICT (codigo_ibge) DO NOTHING
		`)
	case "sqlite3":
		stmt, err = s.db.Prepare(`
			INSERT OR IGNORE INTO cidades (codigo_ibge, nome, codigo_tom, micro_regiao, regiao_imediata, estado_codigo_ibge) 
			VALUES (?, ?, ?, ?, ?, ?)
		`)
	default:
		// MySQL e outros
		stmt, err = s.db.Prepare(`
			INSERT IGNORE INTO cidades (codigo_ibge, nome, codigo_tom, micro_regiao, regiao_imediata, estado_codigo_ibge) 
			VALUES (?, ?, ?, ?, ?, ?)
		`)
	}

	if err != nil {
		return fmt.Errorf("erro ao preparar statement: %w", err)
	}
	defer stmt.Close()

	// Inserir cidades
	count := 0
	for _, cidade := range cidades {
		// Buscar código TOM se existir
		codigoTOM := tomData[cidade.Codigo]
		var codigoTOMPtr *string
		if codigoTOM != "" {
			codigoTOMPtr = &codigoTOM
		}

		if _, err := stmt.Exec(
			cidade.Codigo,
			cidade.Nome,
			codigoTOMPtr,
			fmt.Sprintf("%d", cidade.Microregiao),
			fmt.Sprintf("%d", cidade.RegiaoImediata),
			cidade.CodigoUF,
		); err != nil {
			return fmt.Errorf("erro ao inserir cidade %s (código: %d): %w", cidade.Nome, cidade.Codigo, err)
		}
		count++
	}

	log.Printf("Processadas %d cidades", count)
	return nil
}
