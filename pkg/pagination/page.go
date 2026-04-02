package pagination

type Page struct {
	Page    int
	PerPage int
}

type PerPage struct {
	PerPage int
}

func (p *Page) Init() int {
	if p.Page < 1 {
		p.Page = 1
	}
	return p.Page
}

func (p *PerPage) Init2() int {
	if p.PerPage < 1 || p.PerPage > 100 {
		p.PerPage = 100
	}
	return p.PerPage
}
