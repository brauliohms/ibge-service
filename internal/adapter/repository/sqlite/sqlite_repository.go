package sqlite

import (
	"database/sql"
	"strings"

	"github.com/brauliohms/ibge-service/internal/domain"
	_ "github.com/mattn/go-sqlite3" // Importamos o driver SQLite. O _ significa que o usamos por seus "efeitos colaterais" (registrar-se no pacote database/sql).
	// _ "github.com/glebarez/go-sqlite"
)

// SQLiteRepository é a implementação do repositório que usa um arquivo SQLite como fonte de dados.
type SQLiteRepository struct {
	db *sql.DB
}

// NewSQLiteRepository cria e retorna uma nova instância do repositório SQLite.
func NewSQLiteRepository(filePath string) (*SQLiteRepository, error) {
	// Abre a conexão com o banco de dados. Se o arquivo não existir, ele será criado.
	db, err := sql.Open("sqlite3", filePath)
	if err != nil {
		return nil, err
	}
	// Ping garante que a conexão é válida.
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return &SQLiteRepository{db: db}, nil
}

// FindAllEstados busca todos os estados no banco de dados SQLite.
// Este método é usado para a carga inicial dos dados em memória.
func (r *SQLiteRepository) FindAllEstados() ([]domain.Estado, error) {
	rows, err := r.db.Query("SELECT codigo_ibge, nome, sigla FROM estados ORDER BY nome")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var estados []domain.Estado
	for rows.Next() {
		var e domain.Estado
		if err := rows.Scan(&e.CodigoIBGE, &e.Nome, &e.Sigla); err != nil {
			return nil, err
		}
		estados = append(estados, e)
	}
	return estados, nil
}

// FindAllCidades busca todas as cidades no SQLite e as organiza para a carga inicial.
func (r *SQLiteRepository) FindAllCidades() ([]domain.Cidade, map[string][]domain.Cidade, error) {
	// query := `
	// 	SELECT c.codigo_ibge, c.nome, c.codigo_tom, c.micro_regiao, c.regiao_imediata, e.sigla, e.nome, e.codigo_ibge
	// 	FROM cidades c
	// 	JOIN estados e ON c.estado_codigo_ibge = e.codigo_ibge`
	query := `
		SELECT 
			c.codigo_ibge, 
			c.nome, 
			COALESCE(c.codigo_tom, '') as codigo_tom,
			COALESCE(c.micro_regiao, '') as micro_regiao,
			COALESCE(c.regiao_imediata, '') as regiao_imediata,
			e.sigla,
			e.nome,
			e.codigo_ibge
		FROM cidades c
		INNER JOIN estados e ON c.estado_codigo_ibge = e.codigo_ibge
		ORDER BY e.sigla, c.nome
	`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	var allCidades []domain.Cidade
	cidadesPorEstado := make(map[string][]domain.Cidade)

	for rows.Next() {
		var c domain.Cidade
		// var estadoSigla string
		// Usamos sql.NullString para campos que podem ser nulos, como codigo_tom.
		var codigoTom sql.NullString 

		if err := rows.Scan(&c.CodigoIBGE, &c.Nome, &codigoTom, &c.MicroRegiao, &c.RegiaoImediata, &c.EstadoSigla, &c.EstadoNome, &c.EstadoCodigoIBGE); err != nil {
			return nil, nil, err
		}
		
		if codigoTom.Valid {
			c.CodigoTOM = codigoTom.String
		}

		// Garantimos que a chave do mapa seja sempre maiúscula para consistência.
		ucSigla := strings.ToUpper(c.EstadoSigla)
		allCidades = append(allCidades, c)
		cidadesPorEstado[ucSigla] = append(cidadesPorEstado[ucSigla], c)
	}
	return allCidades, cidadesPorEstado, nil
}

// Close fecha a conexão com o banco de dados.
func (r *SQLiteRepository) Close() {
	r.db.Close()
}
