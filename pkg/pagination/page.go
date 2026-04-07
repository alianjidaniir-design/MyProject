package pagination

import "fmt"

func CheckPage(page, pageSize int) (int, int, error) {
	if page < 1 {
		page = 1
		fmt.Println(page)
	}
	if pageSize < 1 || pageSize > 10 {
		pageSize = 10
	}
	return page, pageSize, nil
}
