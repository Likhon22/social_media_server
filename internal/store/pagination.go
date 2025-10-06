package store

import (
	"net/http"
	"strconv"
)

type PaginatedFeedQuery struct {
	Limit  int    `json:"limit" validate:"gte=1,lte=20"` // gte = greater than or equal
	Page   int    `json:"page" validate:"gte=1"`         // optional
	Offset int    `json:"offset"`
	Sort   string `json:"sort" validate:"oneof=asc desc"`
}

func (fq *PaginatedFeedQuery) Parse(r *http.Request) (PaginatedFeedQuery, error) {
	qs := r.URL.Query()
	limit := qs.Get("limit")
	if limit != "" {
		l, err := strconv.Atoi(limit)
		if err != nil {
			return *fq, err

		}
		fq.Limit = l
	}
	page := qs.Get("page")
	if page != "" {
		p, err := strconv.Atoi(page)
		if err != nil {
			return *fq, err

		}
		fq.Page = p
	}
	fq.Offset = (fq.Page - 1) * fq.Limit
	sort := qs.Get("sort")
	if sort != "" {

		fq.Sort = sort
	}
	return *fq, nil
}
