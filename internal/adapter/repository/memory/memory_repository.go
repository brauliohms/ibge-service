package memory

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/brauliohms/ibge-service/internal/domain"
)

// SourceRepository é uma interface para o repositório que servirá de fonte de dados inicial.
type SourceRepository interface {
	FindAllEstados() ([]domain.Estado, error)
	// Adicionamos um método auxiliar para carregar todas as cidades de forma eficiente.
	FindAllCidades() ([]domain.Cidade, map[string][]domain.Cidade, error)
}

// MemoryRepository implementa a interface IBGERepository e armazena os dados em memória.
type MemoryRepository struct {
	estados           []domain.Estado
	estadosByUF       map[string]domain.Estado
	cidadesByEstadoUF map[string][]domain.Cidade
	cidadesByCodigo   map[string]domain.Cidade // Novo: índice para busca por código
}

// NewMemoryRepository cria e inicializa o repositório em memória, carregando dados da fonte.
func NewMemoryRepository(source SourceRepository) (*MemoryRepository, error) {
	estados, err := source.FindAllEstados()
	if err != nil {
		return nil, fmt.Errorf("falha ao carregar estados: %w", err)
	}

	todasCidades, cidadesMap, err := source.FindAllCidades()
	if err != nil {
		return nil, fmt.Errorf("falha ao carregar cidades: %w", err)
	}

	estadosByUF := make(map[string]domain.Estado)
	for _, e := range estados {
		estadosByUF[strings.ToUpper(e.Sigla)] = e
	}

	// Criar índice de cidades por código IBGE para busca rápida
	cidadesByCodigo := make(map[string]domain.Cidade)
	for _, cidade := range todasCidades {
		cidadesByCodigo[strconv.Itoa(cidade.CodigoIBGE)] = cidade
	}

	return &MemoryRepository{
		estados:           estados,
		estadosByUF:       estadosByUF,
		cidadesByEstadoUF: cidadesMap,
		cidadesByCodigo:   cidadesByCodigo,
	}, nil
}

func (r *MemoryRepository) FindAllEstados() ([]domain.Estado, error) {
	return r.estados, nil
}

func (r *MemoryRepository) FindEstadoByUF(uf string) (*domain.Estado, error) {
	estado, found := r.estadosByUF[strings.ToUpper(uf)]
	if !found {
		return nil, fmt.Errorf("estado com a sigla %s não encontrado", uf)
	}
	return &estado, nil
}

func (r *MemoryRepository) FindCidadesByEstadoUF(uf string) ([]domain.Cidade, error) {
	cidades, found := r.cidadesByEstadoUF[strings.ToUpper(uf)]
	if !found {
		// Verificamos primeiro se o estado existe para dar uma mensagem de erro melhor.
		if _, stateExists := r.estadosByUF[strings.ToUpper(uf)]; !stateExists {
			return nil, fmt.Errorf("estado com a sigla %s não encontrado", uf)
		}
		// O estado existe, mas não tem cidades (cenário improvável, mas possível).
		return []domain.Cidade{}, nil
	}
	return cidades, nil
}

// FindCidadeByCodigo busca uma cidade pelo seu código IBGE
func (r *MemoryRepository) FindCidadeByCodigo(codigo_ibge string) (*domain.Cidade, error) {
	// Validar se o código é um número válido
	if _, err := strconv.Atoi(codigo_ibge); err != nil {
		return nil, fmt.Errorf("código IBGE inválido: %s deve ser um número", codigo_ibge)
	}

	cidade, found := r.cidadesByCodigo[codigo_ibge]
	if !found {
		return nil, fmt.Errorf("cidade com código IBGE %s não encontrada", codigo_ibge)
	}
	
	return &cidade, nil
}
