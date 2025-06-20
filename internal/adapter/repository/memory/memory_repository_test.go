package memory

import (
	"reflect"
	"testing"

	"github.com/brauliohms/ibge-service/internal/domain"
)

// mockSourceRepository é uma implementação falsa da fonte de dados para nossos testes.
type mockSourceRepository struct{}

func (m *mockSourceRepository) FindAllEstados() ([]domain.Estado, error) {
	return []domain.Estado{
		{CodigoIBGE: 1, Nome: "Estado A", Sigla: "EA"},
		{CodigoIBGE: 2, Nome: "Estado B", Sigla: "EB"},
	}, nil
}

func (m *mockSourceRepository) FindAllCidades() ([]domain.Cidade, map[string][]domain.Cidade, error) {
	cidadesMap := map[string][]domain.Cidade{
		"EA": {{CodigoIBGE: 101, Nome: "Cidade A1"}},
		"EB": {{CodigoIBGE: 201, Nome: "Cidade B1"}, {CodigoIBGE: 202, Nome: "Cidade B2"}},
	}
	return nil, cidadesMap, nil // O primeiro retorno não é usado no construtor, podemos deixar nil.
}

func TestMemoryRepository(t *testing.T) {
	// Setup: Criar o repositório em memória usando nosso mock.
	source := &mockSourceRepository{}
	repo, err := NewMemoryRepository(source)
	if err != nil {
		t.Fatalf("Falha ao criar o repositório em memória: %v", err)
	}

	t.Run("deve encontrar um estado pela UF", func(t *testing.T) {
		uf := "EA"
		expected := &domain.Estado{CodigoIBGE: 1, Nome: "Estado A", Sigla: "EA"}

		got, err := repo.FindEstadoByUF(uf)

		if err != nil {
			t.Errorf("Esperava não ter erro, mas recebi: %v", err)
		}
		if !reflect.DeepEqual(got, expected) {
			t.Errorf("Estado retornado incorreto. got: %v, want: %v", got, expected)
		}
	})

	t.Run("deve retornar erro para UF inexistente", func(t *testing.T) {
		uf := "XX"

		_, err := repo.FindEstadoByUF(uf)

		if err == nil {
			t.Errorf("Esperava um erro para UF inexistente, mas não recebi nenhum.")
		}
	})

	t.Run("deve encontrar cidades pela UF do estado", func(t *testing.T) {
		uf := "EB"
		expected := []domain.Cidade{
			{CodigoIBGE: 201, Nome: "Cidade B1"},
			{CodigoIBGE: 202, Nome: "Cidade B2"},
		}

		got, err := repo.FindCidadesByEstadoUF(uf)

		if err != nil {
			t.Errorf("Esperava não ter erro, mas recebi: %v", err)
		}
		if len(got) != len(expected) {
			t.Fatalf("Número de cidades incorreto. got: %d, want: %d", len(got), len(expected))
		}
		if !reflect.DeepEqual(got, expected) {
			t.Errorf("Lista de cidades incorreta. got: %v, want: %v", got, expected)
		}
	})
}
