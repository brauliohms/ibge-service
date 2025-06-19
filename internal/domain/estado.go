package domain

// Estado representa uma Unidade Federativa do Brasil.
type Estado struct {
	CodigoIBGE int    `json:"codigo_ibge"`
	Nome       string `json:"nome"`
	Sigla      string `json:"sigla"`
}
