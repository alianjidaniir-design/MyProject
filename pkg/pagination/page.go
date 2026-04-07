package pagination

type Page struct {
	Page     int
	PageSize int
}

func (p *Page) Init() int {

	if p.Page < 1 {
		p.Page = 1
	}
	if p.PageSize < 1 || p.PageSize > 10 {
		p.PageSize = 10
	}
	offest := (p.Page - 1) * p.PageSize
	return offest
}
