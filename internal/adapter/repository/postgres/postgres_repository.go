package postgres

import (
	"database/sql"
	"strings"

	"github.com/brauliohms/ibge-service/internal/domain"
	_ "github.com/lib/pq" // Driver do PostgreSQL
)

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(dataSourceName string) (*PostgresRepository, error) {
	db, err := sql.Open("postgres", dataSourceName)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return &PostgresRepository{db: db}, nil
}

// FindAllEstados busca todos os estados no banco de dados PostgreSQL.
// Atenção: este repositório é usado apenas para a carga inicial.
func (r *PostgresRepository) FindAllEstados() ([]domain.Estado, error) {
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

// FindEstadoByUF é implementado para satisfazer a interface, mas não será usado em produção (o de memória será).
func (r *PostgresRepository) FindEstadoByUF(uf string) (*domain.Estado, error) {
	// Implementação omitida para brevidade, pois o repositório em memória será o principal.
	return nil, nil
}

// FindCidadesByEstadoUF é implementado para satisfazer a interface, mas não será usado em produção.
func (r *PostgresRepository) FindCidadesByEstadoUF(uf string) ([]domain.Cidade, error) {
	// Implementação omitida para brevidade.
	return nil, nil
}

// FindAllCidades - um método auxiliar para a carga inicial
func (r *PostgresRepository) FindAllCidades() ([]domain.Cidade, map[string][]domain.Cidade, error) {
	// query := `
	// 	SELECT c.codigo_ibge, c.nome, c.codigo_tom, c.micro_regiao, c.regiao_imediata, e.sigla
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

	allCidades := []domain.Cidade{}
	cidadesPorEstado := make(map[string][]domain.Cidade)

	for rows.Next() {
		var c domain.Cidade
		var codigoTom sql.NullString
		// var estadoSigla string
		if err := rows.Scan(&c.CodigoIBGE, &c.Nome, &c.CodigoTOM, &c.MicroRegiao, &c.RegiaoImediata, &c.EstadoSigla, &c.EstadoNome, &c.EstadoCodigoIBGE); err != nil {
			return nil, nil, err
		}
		if codigoTom.Valid {
			c.CodigoTOM = codigoTom.String
		}
		ucSigla := strings.ToUpper(c.EstadoSigla)
		allCidades = append(allCidades, c)
		cidadesPorEstado[ucSigla] = append(cidadesPorEstado[ucSigla], c)
	}
	return allCidades, cidadesPorEstado, nil
}
