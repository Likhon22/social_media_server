package store

import (
	"net/http"
	"strconv"
	"strings"
	"time"
)

type PaginatedFeedQuery struct {
	Limit  int       `json:"limit" validate:"gte=1,lte=20"` // gte = greater than or equal
	Page   int       `json:"page" validate:"gte=1"`         // optional
	Offset int       `json:"offset"`
	Sort   string    `json:"sort" validate:"oneof=asc desc"`
	Tags   []string  `json:"tags" validate:"max=5"`
	Search string    `json:"search" validate:"max=100"`
	Since  time.Time `json:"since"`
	Until  time.Time `json:"until"`
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
	tags := qs.Get("tags")
	if tags != "" {

		fq.Tags = strings.Split(tags, ",")
	}
	search := qs.Get("search")

	if search != "" {

		fq.Search = search
	}
	sinceString := qs.Get("since")
	if sinceString != "" {

		since, err := fq.ParseTime(sinceString)
		if err != nil {
			return *fq, err

		}
		fq.Since = since
	}
	untilStr := qs.Get("until")
	if untilStr != "" {
		until, err := fq.ParseTime(untilStr)
		if err != nil {
			return *fq, err

		}
		fq.Until = until
	}
	return *fq, nil
}

func (fq *PaginatedFeedQuery) ParseTime(timeStr string) (time.Time, error) {
	layout := "2006-01-02T15:04:05Z" // ISO 8601 format
	t, err := time.Parse(layout, timeStr)
	return t, err
}
