package http

import (
	"encoding/json"
	"net/http"
	"strings"

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
		http.Error(w, "Erro interno do servidor", http.StatusInternalServerError)
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
// @Param        uf   path      string  true  "Sigla do Estado (ex: SP, RJ, BA)"
// @Success      200  {object}  domain.Estado "Dados do estado retornados com sucesso"
// @Failure      404  {object}  map[string]string "Estado não encontrado"
// @Router       /estados/{uf} [get]
func (h *IBGEHandler) GetEstadoByUF(w http.ResponseWriter, r *http.Request) {
	uf := chi.URLParam(r, "uf")
	estado, err := h.useCase.GetEstadoByUF(uf)
	if err != nil {
		respondWithJSON(w, http.StatusNotFound, map[string]string{"error": err.Error()})
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
// @Param        uf   path      string  true  "Sigla do Estado (ex: SP, RJ, BA)"
// @Success      200  {array}   domain.Cidade "Lista de cidades retornada com sucesso"
// @Failure      404  {object}  map[string]string "Estado não encontrado"
// @Router       /estados/{uf}/cidades [get]
func (h *IBGEHandler) GetCidadesByEstadoUF(w http.ResponseWriter, r *http.Request) {
	uf := chi.URLParam(r, "uf")
	cidades, err := h.useCase.GetCidadesByEstadoUF(strings.ToUpper(uf))
	if err != nil {
		respondWithJSON(w, http.StatusNotFound, map[string]string{"error": err.Error()})
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
	codigo_ibge := chi.URLParam(r, "codigo_ibge")
	estado, err := h.useCase.GetCidadeByCodigo(codigo_ibge)
	if err != nil {
		respondWithJSON(w, http.StatusNotFound, map[string]string{"error": err.Error()})
		return
	}
	respondWithJSON(w, http.StatusOK, estado)
}

// respondWithJSON é uma função helper para padronizar as respostas JSON.
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
