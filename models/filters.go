package models

import (
	"errors"
	"math"
	"strings"
)

type Filter struct {
	Page     int
	PageSize int
	OrderBy  string
	Query    string
}

type Metadata struct {
	CurrentPage  int
	PageSize     int
	FirstPage    int
	NextPage     int
	PrevPage     int
	LastPage     int
	TotalRecords int
}

func (f *Filter) Validate() error {
	if f.Page <= 0 || f.Page >= 10_000_000 {
		return errors.New("invalid page range: 1 to 10 million")
	}

	if f.PageSize <= 0 || f.PageSize > 100 {
		return errors.New("invalid page size: 1 to 100 max")
	}

	return nil
}
func (f *Filter) addOrdering(q string) string {
	if f.OrderBy == "popular" {
		return strings.Replace(q, "#orderby#", "ORDER BY votes desc, p.created_at desc", 1)
	}

	return strings.Replace(q, "#orderby#", "ORDER BY p.created_at desc", 1)
}

func (f *Filter) addWhere(q string) string {
	if len(f.Query) > 0 {
		return strings.Replace(q, "#where#", "WHERE LOWER(p.title) LIKE $1", 1)
	}

	return strings.Replace(q, "#where#", "", 1)
}

func (f *Filter) addLimitOffset(q string) string {
	if len(f.Query) > 0 {
		return strings.Replace(q, "#limit#", "LIMIT $2 OFFSET $3", 1)
	}

	return strings.Replace(q, "#limit#", "LIMIT $1 OFFSET $2", 1)
}

func (f *Filter) applyTemplate(q string) string {
	return f.addLimitOffset(f.addWhere(f.addOrdering(q)))
}

func (f *Filter) limit() int {
	return f.PageSize
}

func (f *Filter) offset() int {
	return (f.Page - 1) * f.PageSize
}

func calculateMetadata(totalRecords, page, pageSize int) Metadata {
	if totalRecords == 0 {
		return Metadata{}
	}

	meta := Metadata{
		CurrentPage:  page,
		PageSize:     pageSize,
		FirstPage:    1,
		LastPage:     int(math.Ceil(float64(totalRecords) / float64(pageSize))),
		TotalRecords: totalRecords,
	}

	meta.NextPage = meta.CurrentPage + 1
	meta.PrevPage = meta.CurrentPage - 1

	if meta.CurrentPage <= meta.FirstPage {
		meta.PrevPage = 0
	}

	return meta
}
