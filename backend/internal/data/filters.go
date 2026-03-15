package data

import (
	"slices"
	"strings"

	"github.com/OTDG-Dev/On-Track-Hockey/backend/internal/validator"
)

type Filters struct {
	Page         int
	PageSize     int
	Sort         string
	SortSafeList []string // hold supported sort values
}

func ValidateFilters(v *validator.Validator, f Filters) {
	v.Check(f.Page > 0, "page", "must be greater than zero")
	v.Check(f.Page <= 10_000_000, "page", "must be a maximum of 10 million")
	v.Check(f.PageSize > 0, "page_size", "must be greater than zero")
	v.Check(f.PageSize <= 100, "page_size", "must have a maximum of 100")

	v.Check(validator.PermittedValue(f.Sort, f.SortSafeList...), "sort", "invalid sort value")
}

// check that the client-provided Sort field matches an entry on the safe list (whitelist)
// if so get the column name striping leading - if one exists
func (f Filters) sortColumn() string {
	if slices.Contains(f.SortSafeList, f.Sort) {
		return strings.TrimPrefix(f.Sort, "-")
	}

	// sensible as failsafe for sql injection (we already checked in ValidateFilters)
	panic("unsafe sort paramter: " + f.Sort)
}

// strip the "-" if it exists and convert to sql ordering
func (f Filters) SortDirection() string {
	if strings.HasPrefix(f.Sort, "-") {
		return "DESC"
	}
	return "ASC"
}

func (f Filters) limit() int {
	return f.PageSize
}

// could pontentially cause an overflow, however we limited in validation checks
func (f Filters) offset() int {
	return (f.Page - 1) * f.PageSize
}

type Metadata struct {
	CurrentPage  int `json:"current_page,omitzero"`
	PageSize     int `json:"page_size,omitzero"`
	FirstPage    int `json:"first_page,omitzero"`
	LastPage     int `json:"last_page,omitzero"`
	TotalRecords int `json:"total_records,omitzero"`
}

// calculates the appropriate pagination metadata values given the total number of records
func calculateMetadata(TotalRecords, page, pageSize int) Metadata {
	if TotalRecords == 0 {
		return Metadata{}
	}

	return Metadata{
		CurrentPage:  page,
		PageSize:     pageSize,
		FirstPage:    1,
		LastPage:     (TotalRecords + pageSize - 1) / pageSize,
		TotalRecords: TotalRecords,
	}
}
