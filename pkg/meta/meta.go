package meta

type Meta struct {
	Page       int32 `json:"page"`
	PerPage    int32 `json:"per_page"`
	PageCount  int32 `json:"page_count"`
	TotalCount int32 `json:"total_count"`
}

func New(page, perPage, total, pagLimitDef int32) (*Meta, error) {
	if perPage <= 0 {
		perPage = pagLimitDef
	}

	pageCount := int32(0)
	if total >= 0 {
		pageCount = (total + perPage - 1) / perPage
		if page > pageCount {
			page = pageCount
		}
	}

	if page < 1 {
		page = 1
	}

	return &Meta{
		Page:       page,
		PerPage:    perPage,
		PageCount:  pageCount,
		TotalCount: total,
	}, nil
}

func (meta *Meta) Offset() int32 {
	return (meta.Page - 1) * meta.PerPage
}

func (meta *Meta) Limit() int32 {
	return meta.PerPage
}