package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/brauliohms/ibge-service/internal/domain"
	"github.com/brauliohms/ibge-service/internal/usecase"
)

// mockIBGERepository é um mock para o repositório usado pelo caso de uso.
type mockIBGERepository struct{}

func (m *mockIBGERepository) FindAllEstados() ([]domain.Estado, error) {
	return []domain.Estado{
		{CodigoIBGE: 35, Sigla: "SP", Nome: "São Paulo"},
		{CodigoIBGE: 33, Sigla: "RJ", Nome: "Rio de Janeiro"},
		{CodigoIBGE: 31, Sigla: "MG", Nome: "Minas Gerais"},
	}, nil
}

func (m *mockIBGERepository) FindEstadoByUF(uf string) (*domain.Estado, error) {
	uf = strings.ToUpper(uf)
	switch uf {
	case "SP":
		return &domain.Estado{CodigoIBGE: 35, Sigla: "SP", Nome: "São Paulo"}, nil
	case "RJ":
		return &domain.Estado{CodigoIBGE: 33, Sigla: "RJ", Nome: "Rio de Janeiro"}, nil
	case "MG":
		return &domain.Estado{CodigoIBGE: 31, Sigla: "MG", Nome: "Minas Gerais"}, nil
	default:
		return nil, fmt.Errorf("estado com a sigla %s não encontrado", uf)
	}
}

func (m *mockIBGERepository) FindCidadesByEstadoUF(uf string) ([]domain.Cidade, error) {
	uf = strings.ToUpper(uf)
	switch uf {
	case "SP":
		return []domain.Cidade{
			{CodigoIBGE: 3550308, Nome: "São Paulo", EstadoCodigoIBGE: 35},
			{CodigoIBGE: 3509502, Nome: "Campinas", EstadoCodigoIBGE: 35},
			{CodigoIBGE: 3552205, Nome: "Santos", EstadoCodigoIBGE: 35},
		}, nil
	case "RJ":
		return []domain.Cidade{
			{CodigoIBGE: 3304557, Nome: "Rio de Janeiro", EstadoCodigoIBGE: 33},
			{CodigoIBGE: 3301702, Nome: "Niterói", EstadoCodigoIBGE: 33},
		}, nil
	case "MG":
		return []domain.Cidade{
			{CodigoIBGE: 3106200, Nome: "Belo Horizonte", EstadoCodigoIBGE: 31},
		}, nil
	default:
		return nil, fmt.Errorf("estado com a sigla %s não encontrado", uf)
	}
}

func TestIBGEHandler(t *testing.T) {
	// Setup: criar as camadas com o mock
	repo := &mockIBGERepository{}
	uc := usecase.NewIBGEUseCase(repo)
	handler := NewIBGEHandler(uc)
	router := SetupRouter(handler)

	t.Run("GET /api/v1/estados - deve retornar todos os estados", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/estados", nil)
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		// Verificar status code
		if status := rr.Code; status != http.StatusOK {
			t.Errorf("Status code incorreto: got %v want %v", status, http.StatusOK)
		}

		// Verificar Content-Type
		expectedContentType := "application/json"
		if contentType := rr.Header().Get("Content-Type"); contentType != expectedContentType {
			t.Errorf("Content-Type incorreto: got %v want %v", contentType, expectedContentType)
		}

		// Verificar resposta JSON
		var estados []domain.Estado
		if err := json.Unmarshal(rr.Body.Bytes(), &estados); err != nil {
			t.Fatalf("Erro ao decodificar JSON: %v", err)
		}

		if len(estados) != 3 {
			t.Errorf("Número de estados incorreto: got %d want %d", len(estados), 3)
		}

		// Verificar se contém os estados esperados
		estadosMap := make(map[string]domain.Estado)
		for _, estado := range estados {
			estadosMap[estado.Sigla] = estado
		}

		expectedEstados := []string{"SP", "RJ", "MG"}
		for _, sigla := range expectedEstados {
			if _, exists := estadosMap[sigla]; !exists {
				t.Errorf("Estado %s não encontrado na resposta", sigla)
			}
		}
	})

	t.Run("GET /api/v1/estados/{uf} - deve retornar estado específico", func(t *testing.T) {
		testCases := []struct {
			uf           string
			expectedCode int
			expectedName string
		}{
			{"SP", http.StatusOK, "São Paulo"},
			{"sp", http.StatusOK, "São Paulo"}, // Teste case insensitive
			{"RJ", http.StatusOK, "Rio de Janeiro"},
			{"MG", http.StatusOK, "Minas Gerais"},
		}

		for _, tc := range testCases {
			t.Run("UF_"+tc.uf, func(t *testing.T) {
				req := httptest.NewRequest("GET", "/api/v1/estados/"+tc.uf, nil)
				rr := httptest.NewRecorder()

				router.ServeHTTP(rr, req)

				if status := rr.Code; status != tc.expectedCode {
					t.Errorf("Status code incorreto para %s: got %v want %v", tc.uf, status, tc.expectedCode)
				}

				if tc.expectedCode == http.StatusOK {
					var estado domain.Estado
					if err := json.Unmarshal(rr.Body.Bytes(), &estado); err != nil {
						t.Fatalf("Erro ao decodificar JSON: %v", err)
					}

					if estado.Nome != tc.expectedName {
						t.Errorf("Nome do estado incorreto: got %v want %v", estado.Nome, tc.expectedName)
					}

					if estado.Sigla != strings.ToUpper(tc.uf) {
						t.Errorf("Sigla do estado incorreta: got %v want %v", estado.Sigla, strings.ToUpper(tc.uf))
					}
				}
			})
		}
	})

	t.Run("GET /api/v1/estados/{uf} - deve retornar 404 para estado inexistente", func(t *testing.T) {
		invalidUFs := []string{"XX", "YY", "ZZ", "ABC"}

		for _, uf := range invalidUFs {
			t.Run("Invalid_UF_"+uf, func(t *testing.T) {
				req := httptest.NewRequest("GET", "/api/v1/estados/"+uf, nil)
				rr := httptest.NewRecorder()

				router.ServeHTTP(rr, req)

				if status := rr.Code; status != http.StatusNotFound {
					t.Errorf("Status code incorreto para UF inválida %s: got %v want %v", uf, status, http.StatusNotFound)
				}
			})
		}
	})

	t.Run("GET /api/v1/estados/{uf}/cidades - deve retornar cidades do estado", func(t *testing.T) {
		testCases := []struct {
			uf            string
			expectedCount int
			expectedCities []string
		}{
			{"SP", 3, []string{"São Paulo", "Campinas", "Santos"}},
			{"RJ", 2, []string{"Rio de Janeiro", "Niterói"}},
			{"MG", 1, []string{"Belo Horizonte"}},
		}

		for _, tc := range testCases {
			t.Run("Cidades_"+tc.uf, func(t *testing.T) {
				req := httptest.NewRequest("GET", "/api/v1/estados/"+tc.uf+"/cidades", nil)
				rr := httptest.NewRecorder()

				router.ServeHTTP(rr, req)

				if status := rr.Code; status != http.StatusOK {
					t.Errorf("Status code incorreto: got %v want %v", status, http.StatusOK)
				}

				var cidades []domain.Cidade
				if err := json.Unmarshal(rr.Body.Bytes(), &cidades); err != nil {
					t.Fatalf("Erro ao decodificar JSON: %v", err)
				}

				if len(cidades) != tc.expectedCount {
					t.Errorf("Número de cidades incorreto para %s: got %d want %d", tc.uf, len(cidades), tc.expectedCount)
				}

				// Verificar se contém as cidades esperadas
				cidadesMap := make(map[string]bool)
				for _, cidade := range cidades {
					cidadesMap[cidade.Nome] = true
				}

				for _, expectedCity := range tc.expectedCities {
					if !cidadesMap[expectedCity] {
						t.Errorf("Cidade %s não encontrada para estado %s", expectedCity, tc.uf)
					}
				}
			})
		}
	})

	t.Run("GET /api/v1/estados/{uf}/cidades - deve retornar 404 para estado inexistente", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/estados/XX/cidades", nil)
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusNotFound {
			t.Errorf("Status code incorreto: got %v want %v", status, http.StatusNotFound)
		}
	})

	t.Run("GET /api/v1/cidades/{codigo} - deve retornar cidade específica", func(t *testing.T) {
		testCases := []struct {
			codigo       string
			expectedName string
			expectedUF   int
		}{
			{"3550308", "São Paulo", 35},
			{"3509502", "Campinas", 35},
			{"3304557", "Rio de Janeiro", 33},
			{"3106200", "Belo Horizonte", 31},
		}

		for _, tc := range testCases {
			t.Run("Codigo_"+tc.codigo, func(t *testing.T) {
				req := httptest.NewRequest("GET", "/api/v1/cidades/"+tc.codigo, nil)
				rr := httptest.NewRecorder()

				router.ServeHTTP(rr, req)

				if status := rr.Code; status != http.StatusOK {
					t.Errorf("Status code incorreto: got %v want %v", status, http.StatusOK)
				}

				var cidade domain.Cidade
				if err := json.Unmarshal(rr.Body.Bytes(), &cidade); err != nil {
					t.Fatalf("Erro ao decodificar JSON: %v", err)
				}

				if cidade.Nome != tc.expectedName {
					t.Errorf("Nome da cidade incorreto: got %v want %v", cidade.Nome, tc.expectedName)
				}

				if cidade.EstadoCodigoIBGE != tc.expectedUF {
					t.Errorf("Código do estado incorreto: got %v want %v", cidade.EstadoCodigoIBGE, tc.expectedUF)
				}
			})
		}
	})

	t.Run("GET /api/v1/cidades/{codigo} - deve retornar 404 para cidade inexistente", func(t *testing.T) {
		invalidCodes := []string{"999999", "000000", "123456"}

		for _, code := range invalidCodes {
			t.Run("Invalid_Code_"+code, func(t *testing.T) {
				req := httptest.NewRequest("GET", "/api/v1/cidades/"+code, nil)
				rr := httptest.NewRecorder()

				router.ServeHTTP(rr, req)

				if status := rr.Code; status != http.StatusNotFound {
					t.Errorf("Status code incorreto para código inválido %s: got %v want %v", code, status, http.StatusNotFound)
				}
			})
		}
	})
}

// Teste de benchmark para verificar performance
func BenchmarkGetAllEstados(b *testing.B) {
	repo := &mockIBGERepository{}
	uc := usecase.NewIBGEUseCase(repo)
	handler := NewIBGEHandler(uc)
	router := SetupRouter(handler)

	req := httptest.NewRequest("GET", "/api/v1/estados", nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)
	}
}

// Teste de benchmark para busca de cidades por estado
func BenchmarkGetCidadesByEstado(b *testing.B) {
	repo := &mockIBGERepository{}
	uc := usecase.NewIBGEUseCase(repo)
	handler := NewIBGEHandler(uc)
	router := SetupRouter(handler)

	req := httptest.NewRequest("GET", "/api/v1/estados/SP/cidades", nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)
	}
}

// Teste para verificar headers de resposta
func TestResponseHeaders(t *testing.T) {
	repo := &mockIBGERepository{}
	uc := usecase.NewIBGEUseCase(repo)
	handler := NewIBGEHandler(uc)
	router := SetupRouter(handler)

	req := httptest.NewRequest("GET", "/api/v1/estados", nil)
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	// Verificar Content-Type
	contentType := rr.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Content-Type incorreto: got %v want %v", contentType, "application/json")
	}

	// Verificar headers de segurança (se implementados)
	xContentTypeOptions := rr.Header().Get("X-Content-Type-Options")
	if xContentTypeOptions != "nosniff" {
		t.Errorf("X-Content-Type-Options incorreto: got %v want %v", xContentTypeOptions, "nosniff")
	}
}

// Teste para verificar rate limiting (se implementado)
func TestRateLimiting(t *testing.T) {
	repo := &mockIBGERepository{}
	uc := usecase.NewIBGEUseCase(repo)
	handler := NewIBGEHandler(uc)
	router := SetupRouter(handler)

	// Fazer várias requisições rapidamente
	for i := 0; i < 10; i++ {
		req := httptest.NewRequest("GET", "/api/v1/estados", nil)
		req.RemoteAddr = "192.168.1.1:12345" // Simular mesmo IP
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		// As primeiras requisições devem passar
		if i < 5 && rr.Code != http.StatusOK {
			t.Errorf("Requisição %d deveria passar: got status %d", i, rr.Code)
		}
	}
}

func (m *mockIBGERepository) FindCidadeByCodigo(codigo string) (*domain.Cidade, error) {
	switch codigo {
	case "3550308":
		return &domain.Cidade{CodigoIBGE: 3550308, Nome: "São Paulo", EstadoCodigoIBGE: 35}, nil
	case "3509502":
		return &domain.Cidade{CodigoIBGE: 3509502, Nome: "Campinas", EstadoCodigoIBGE: 35}, nil
	case "3552205":
		return &domain.Cidade{CodigoIBGE: 3552205, Nome: "Santos", EstadoCodigoIBGE: 35}, nil
	case "3304557":
		return &domain.Cidade{CodigoIBGE: 3304557, Nome: "Rio de Janeiro", EstadoCodigoIBGE: 33}, nil
	case "3301702":
		return &domain.Cidade{CodigoIBGE: 3301702, Nome: "Niterói", EstadoCodigoIBGE: 33}, nil
	case "3106200":
		return &domain.Cidade{CodigoIBGE: 3106200, Nome: "Belo Horizonte", EstadoCodigoIBGE: 31}, nil
	default:
		return nil, fmt.Errorf("cidade com código %s não encontrada", codigo)
	}
}