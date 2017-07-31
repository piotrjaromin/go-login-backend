package common

type Order string

var (
	Desc Order = "desc"
	Asc Order = "asc"
)

type SortField struct {
	Name string
	Order Order
}

type SortFields []SortField