package httputil

import (
	"net/http"
	"strconv"
)

const (
	PaginationMaxSize = 50
)

type Pager struct {
	Page   int
	Size   int
	OffSet int
}

//GetPager returns pager object containing pagination params
func GetPager(r *http.Request) Pager {

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	size, _ := strconv.Atoi(r.URL.Query().Get("size"))

	maxSize := PaginationMaxSize
	if size > maxSize || size <= 0 {
		size = maxSize
	}

	if page <= 0 {
		page = 1
	}

	offset := (page - 1) * size

	return Pager{
		Page:   page,
		Size:   size,
		OffSet: offset,
	}
}
