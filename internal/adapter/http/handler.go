package http

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/brauliohms/ibge-service/internal/domain"
	"github.com/brauliohms/ibge-service/internal/usecase"
	"github.com/go-chi/chi/v5"
)

type IBGEHandler struct {
	useCase *usecase.IBGEUseCase
}

func NewIBGEHandler(uc *usecase.IBGEUseCase) *IBGEHandler {
	return &IBGEHandler{useCase: uc}
}

// GetAllEstados godoc
// @Summary      Lista todos os estados
// @Description  Retorna um array com todos os 27 estados brasileiros
// @Tags         Estados
// @Accept       json
// @Produce      json
// @Success      200  {array}   domain.Estado "Lista de estados retornada com sucesso"
// @Failure      500  {object}  map[string]string "Erro interno do servidor"
// @Router       /estados [get]
func (h *IBGEHandler) GetAllEstados(w http.ResponseWriter, r *http.Request) {
	estados, err := h.useCase.GetAllEstados()
	if err != nil {
		log.Printf("Erro ao buscar estados: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Erro interno do servidor")
		return
	}
	respondWithJSON(w, http.StatusOK, estados)
}

// GetEstadoByUF godoc
// @Summary      Busca um estado pela sua sigla (UF)
// @Description  Retorna os dados completos de um único estado
// @Tags         Estados
// @Accept       json
// @Produce      json
// @Param        uf   path      string  true  "Sigla do Estado (ex: SP, RJ, BA) ou Código IBGE do Estado (ex: 35, 33, 29)"
// @Success      200  {object}  domain.Estado "Dados do estado retornados com sucesso"
// @Failure      404  {object}  map[string]string "Estado não encontrado"
// @Router       /estados/{uf} [get]
func (h *IBGEHandler) GetEstadoByUF(w http.ResponseWriter, r *http.Request) {
	ufOuCodigo := chi.URLParam(r, "uf")
	var estado *domain.Estado
	var err error

	// Verifica se é um número (código IBGE) ou string (sigla)
	if _, parseErr := strconv.Atoi(ufOuCodigo); parseErr == nil {
		// É um número, busca por código IBGE
		estado, err = h.useCase.GetEstadoByCodigoIbge(ufOuCodigo)
	} else {
		// É uma string, busca por sigla
		estado, err = h.useCase.GetEstadoByUF(strings.ToUpper(ufOuCodigo))
	}

	if err != nil {
		log.Printf("Erro ao buscar estado %s: %v", ufOuCodigo, err)
		respondWithError(w, http.StatusNotFound, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, estado)
}

// GetCidadesByEstadoUF godoc
// @Summary      Busca todas as cidades de um estado
// @Description  Retorna um array com todas as cidades pertencentes a um determinado estado (UF)
// @Tags         Cidades
// @Accept       json
// @Produce      json
// @Param        uf   path      string  true  "Sigla do Estado (ex: SP, RJ, BA) ou Código IBGE do Estado (ex: 35, 33, 29)"
// @Success      200  {array}   domain.Cidade "Lista de cidades retornada com sucesso"
// @Failure      404  {object}  map[string]string "Estado não encontrado"
// @Router       /estados/{uf}/cidades [get]
func (h *IBGEHandler) GetCidadesByEstadoUF(w http.ResponseWriter, r *http.Request) {
	ufOuCodigo := chi.URLParam(r, "uf")
	var cidades []domain.Cidade
	var err error

	// Verifica se é um número (código IBGE) ou string (sigla)
	if _, parseErr := strconv.Atoi(ufOuCodigo); parseErr == nil {
		// É um número, busca por código IBGE
		cidades, err = h.useCase.GetCidadesByEstadoCodigoIbge(ufOuCodigo)
	} else {
		// É uma string, busca por sigla
		cidades, err = h.useCase.GetCidadesByEstadoUF(strings.ToUpper(ufOuCodigo))
	}

	if err != nil {
		log.Printf("Erro ao buscar cidades do estado %s: %v", ufOuCodigo, err)
		respondWithError(w, http.StatusNotFound, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, cidades)
}

// GetCidadeByCodigo godoc
// @Summary Busca cidade por código IBGE
// @Description Retorna uma cidade específica pelo seu código IBGE
// @Tags Cidades
// @Accept json
// @Produce json
// @Param codigo_ibge path string true "Código IBGE da cidade" example(3550308)
// @Success 200 {object} domain.Cidade
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /cidades/{codigo_ibge} [get]
func (h *IBGEHandler) GetCidadeByCodigo(w http.ResponseWriter, r *http.Request) {
	codigoIBGE := chi.URLParam(r, "codigo_ibge")
	cidade, err := h.useCase.GetCidadeByCodigo(codigoIBGE)
	if err != nil {
		log.Printf("Erro ao buscar cidade com código %s: %v", codigoIBGE, err)
		respondWithError(w, http.StatusNotFound, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, cidade)
}

// GetCidadeByCodigoTOM godoc
// @Summary Busca cidade por código TOM
// @Description Retorna uma cidade específica pelo seu código TOM
// @Tags Cidades
// @Accept json
// @Produce json
// @Param codigo_tom path string true "Código TOM da cidade" example(7107)
// @Success 200 {object} domain.Cidade
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /cidades/{codigo_tom}/tom [get]
func (h *IBGEHandler) GetCidadeByCodigoTOM(w http.ResponseWriter, r *http.Request) {
	codigoTOM := chi.URLParam(r, "codigo_tom")
	cidade, err := h.useCase.GetCidadeByCodigoTOM(codigoTOM)
	if err != nil {
		log.Printf("Erro ao buscar cidade com código TOM %s: %v", codigoTOM, err)
		respondWithError(w, http.StatusNotFound, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, cidade)
}

// respondWithJSON é uma função helper para padronizar as respostas JSON.
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Erro ao fazer marshal do JSON: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Erro interno do servidor")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	// Ignoramos o erro de w.Write intencionalmente aqui
	// pois não há muito que possamos fazer se falhar
	_, _ = w.Write(response)
}

// respondWithError é uma função helper para padronizar as respostas de erro.
func respondWithError(w http.ResponseWriter, code int, message string) {
	errorResponse := map[string]string{"error": message}
	response, err := json.Marshal(errorResponse)
	if err != nil {
		log.Printf("Erro ao fazer marshal do erro: %v", err)
		http.Error(w, "Erro interno do servidor", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	// Ignoramos o erro de w.Write intencionalmente aqui
	// pois não há muito que possamos fazer se falhar
	_, _ = w.Write(response)
}
