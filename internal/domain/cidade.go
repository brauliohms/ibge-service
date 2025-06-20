package domain

// Cidade representa um município brasileiro.
type Cidade struct {
	CodigoIBGE       int    `json:"codigo_ibge"`
	Nome             string `json:"nome"`
	CodigoTOM        string `json:"codigo_tom,omitempty"`
	MicroRegiao      string `json:"micro_regiao,omitempty"`
	RegiaoImediata   string `json:"regiao_imediata,omitempty"`
	EstadoCodigoIBGE int    `json:"estado_codigo_ibge"`
	EstadoSigla      string `json:"estado_sigla"`
	EstadoNome       string `json:"estado_nome"`
}
