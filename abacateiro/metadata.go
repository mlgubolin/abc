package application

import (
	"math"
	"strings"
)

var (
	Environment string
	BuildTime   string
	Version     string
)

const (
	MaxLabelKeyLen = 100
	MaxLabelValLen = 100
)

// Metadata contém informações sobre a paginação e o total de registros retornados.
type Metadata struct {
	CurrentPage  int `json:"current_page,omitempty"`
	PageSize     int `json:"page_size,omitempty"`
	FirstPage    int `json:"first_page,omitempty"`
	LastPage     int `json:"last_page,omitempty"`
	TotalRecords int `json:"total_records,omitempty"`
	TotalCount   int `json:"total_count,omitempty"`
}

// CalculateMetadata calcula as informações de Metadata com base nos parâmetros fornecidos.
func CalculateMetadata(totalRecords, page, pageSize int) Metadata {
	if totalRecords == 0 {
		return Metadata{
			CurrentPage:  page,
			PageSize:     pageSize,
			FirstPage:    1,
			LastPage:     1,
			TotalRecords: 0,
		}
	}

	return Metadata{
		CurrentPage:  page,
		PageSize:     pageSize,
		FirstPage:    1,
		LastPage:     int(math.Ceil(float64(totalRecords) / float64(pageSize))),
		TotalRecords: totalRecords,
	}
}

// Pagination contém informações sobre a paginação e a ordenação dos resultados.
type Pagination struct {
	Sort              string   `json:"sort"`
	SortBy            string   `json:"sort_by"`
	SortSafeList      []string `json:"-"`
	Page              int      `json:"page"`
	PageSize          int      `json:"page_size"`
	SortDescending    bool     `json:"sort_descending"`
	DisablePagination bool     `json:"disable_pagination"`
}

// Limit retorna o valor de limite (quantidade de resultados por página).
func (p *Pagination) Limit() int {
	return p.PageSize
}

func (p *Pagination) Offset() int {
	return (p.Page - 1) * p.PageSize
}

// SortColumn retorna a coluna de ordenação considerando os valoers permitidos em SortSafeList.
// Se a coluna de ordenação estiver no formato "-column", o prefixo "-" será removido antes de retornar o nome da coluna.
// Caso nenhum valor seja encontrado retorna o primeiro valor permitido em SortSafeList ou retorna uma string vazia.
// Caso nenhum valor seja encontrado retorna o primeiro valor permitido em SortSafeList ou retorna uma string vazia.
func (p *Pagination) SortColumn() string {
	for _, safeValue := range p.SortSafeList {
		if p.Sort == safeValue {
			return strings.TrimPrefix(p.Sort, "-")
		}
	}
	if len(p.SortSafeList) > 0 {
		return strings.TrimPrefix(p.SortSafeList[0], "-")
	}
	return ""
}

// SortDirection retorna a direção de ordenação com base na coluna de ordenação.
// Se a coluna de ordenação começar com o prefixo "-", retorna "DESC" (ordem descendente).
// Caso contrário, retorna "ASC" (ordem ascendente).
func (p *Pagination) SortDirection() string {
	if strings.HasPrefix(p.Sort, "-") {
		return "DESC"
	}
	return "ASC"
}

// LimitPagination define os limites padrão para a paginação caso estejam fora dos limites estipulados.
func (p *Pagination) LimitPagination() {
	if p.Page <= 0 {
		p.Page = 1
	}

	if p.PageSize == 0 {
		p.PageSize = 1
	}

	if p.PageSize > 1000 {
		p.PageSize = 1000
	}
}
