package pagination

type Page struct {
	Page    int
	PerPage int
}

type PerPage struct {
	PerPage int
}

func (p *Page) Init(page, perPage int) {
	p.Page = page
	p.PerPage = perPage
	if p.Page < 1 {
		p.Page = 1
	}
	if p.PerPage < 1 || p.PerPage > 100 {
		p.PerPage = 100
	}

}
