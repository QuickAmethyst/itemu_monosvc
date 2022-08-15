package sql

type Paging struct {
	CurrentPage uint `json:"currentPage"`
	PageSize    uint `json:"pageSize"`
	Total       uint `json:"total"`
}

// NewPaging return a Paging instance with default value.
func NewPaging(total uint) *Paging {
	return &Paging{
		CurrentPage: 1,
		PageSize:    12,
		Total:       total,
	}
}

func (p *Paging) BuildQuery() (limitClause string, limitClauseArgs []interface{}) {
	limitArg := p.PageSize
	offsetArg := p.PageSize * (p.CurrentPage - 1)
	limitClause = "LIMIT ? OFFSET ?"
	limitClauseArgs = []interface{}{limitArg, offsetArg}
	return
}

func (p *Paging) LastPage() uint {
	totalPage := p.Total / p.PageSize
	if hasRemaining := p.Total%p.PageSize != 0; hasRemaining {
		totalPage += 1
	}

	return totalPage
}

func (p *Paging) Normalize() {
	if p.PageSize == 0 {
		p.PageSize = 12
	}

	if p.CurrentPage == 0 {
		p.CurrentPage = 1
	}
}
