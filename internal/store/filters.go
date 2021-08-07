package store

import "strings"

type Filters struct {
	Page  int    `json:"page" validate:"gt=0"`
	Limit int    `json:"limit" validate:"gt=0,lt=100"`
	Sort  string `json:"sort" validate:"oneof='id' 'title' '-id' '-title'"`
}

func (f Filters) sortColumn() string {
	return strings.TrimPrefix(f.Sort, "-")
}

func (f Filters) sortDirection() string {
	if strings.HasPrefix(f.Sort, "-") {
		return "DESC"
	}

	return "ASC"
}