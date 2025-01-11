package utils

import (
	"math"

	"gorm.io/gorm"
)

type (
	paginate struct {
		limit int
		total int
		page  int
	}

	Paginator interface {
		Page() int
		PageNums() int
		HasNext() bool
		HasPrevious() bool
		NextPage() int
		PreviousPage() int
	}
)

func NewPaginator(limit, page, total int) Paginator {
	return &paginate{limit: limit, page: page}
}

func (p *paginate) PaginatedResult(db *gorm.DB) *gorm.DB {
	offset := (p.page - 1) * p.limit

	return db.Offset(offset).
		Limit(p.limit)
}

func (p *paginate) Page() int {
	return p.page
}

func (p *paginate) PageNums() int {
	n := int64(math.Ceil(float64(p.total) / float64(p.limit)))
	return int(n)
}

func (p *paginate) HasNext() bool {
	return p.page < p.PageNums()
}

func (p *paginate) HasPrevious() bool {
	return p.page > 1
}

func (p *paginate) NextPage() int {
	return p.page + 1
}

func (p *paginate) PreviousPage() int {
	return p.page - 1
}
