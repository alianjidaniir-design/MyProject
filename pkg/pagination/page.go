package pagination

func CheckPage(page, pageSize int) (int, int) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 10 {
		pageSize = 10
	}
	return page, pageSize
}
