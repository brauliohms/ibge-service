package usecase

import "github.com/brauliohms/ibge-service/internal/domain"

// IBGERepository é a interface que define os contratos de acesso aos dados.
// É a porta de entrada para a persistência, permitindo a inversão de dependência.
type IBGERepository interface {
	FindAllEstados() ([]domain.Estado, error)
	FindEstadoByUF(uf string) (*domain.Estado, error)
	FindEstadoByCodigoIbge(codigo_ibge string) (*domain.Estado, error)
	FindCidadesByEstadoUF(uf string) ([]domain.Cidade, error)
	FindCidadesByEstadoCodigoIbge(codigo_ibge string) ([]domain.Cidade, error)
	FindCidadeByCodigo(codigo_ibge string) (*domain.Cidade, error)
	FindCidadeByCodigoTOM(codigo_tom string) (*domain.Cidade, error)
}

// IBGEUseCase encapsula a lógica de negócio relacionada ao IBGE.
type IBGEUseCase struct {
	repo IBGERepository
}

// NewIBGEUseCase cria uma nova instância do caso de uso, injetando o repositório.
func NewIBGEUseCase(repo IBGERepository) *IBGEUseCase {
	return &IBGEUseCase{repo: repo}
}

// GetAllEstados retorna todos os estados.
func (uc *IBGEUseCase) GetAllEstados() ([]domain.Estado, error) {
	return uc.repo.FindAllEstados()
}

// GetEstadoByUF retorna um estado pela sua sigla.
func (uc *IBGEUseCase) GetEstadoByUF(uf string) (*domain.Estado, error) {
	return uc.repo.FindEstadoByUF(uf)
}

// GetEstadoByCodigoIbge retorna um estado pela seu código IBGE.
func (uc *IBGEUseCase) GetEstadoByCodigoIbge(codigo_ibge string) (*domain.Estado, error) {
	return uc.repo.FindEstadoByCodigoIbge(codigo_ibge)
}

// GetCidadesByEstadoUF retorna todas as cidades de um estado.
func (uc *IBGEUseCase) GetCidadesByEstadoUF(uf string) ([]domain.Cidade, error) {
	return uc.repo.FindCidadesByEstadoUF(uf)
}

// GetCidadesByEstadoCodigoIbge retorna todas as cidades de um estado.
func (uc *IBGEUseCase) GetCidadesByEstadoCodigoIbge(codigo_ibge string) ([]domain.Cidade, error) {
	return uc.repo.FindCidadesByEstadoCodigoIbge(codigo_ibge)
}

func (uc *IBGEUseCase) GetCidadeByCodigo(codigo_ibge string) (*domain.Cidade, error) {
	return uc.repo.FindCidadeByCodigo(codigo_ibge)
}

func (uc *IBGEUseCase) GetCidadeByCodigoTOM(codigo_tom string) (*domain.Cidade, error) {
	return uc.repo.FindCidadeByCodigoTOM(codigo_tom)
}
