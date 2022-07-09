package okresponse

// common response structs:
type (

	// Pager :
	Pager struct {
		Total    int64 `json:"total"`    // 数据总数
		Current  int64 `json:"current"`  // 当前页
		PageSize int64 `json:"pageSize"` // 页大小
		MaxPage  int64 `json:"maxPage"`  // 最大页数
	}
)

////分页
func ToPager(count int64, currentPage int64, pageSize int64) *Pager {
	// page obj:
	pager := &Pager{}
	pager.Total = count

	// page size:
	pager.PageSize = pageSize
	if pageSize < 1 || pageSize > 10000 {
		pager.PageSize = 20
	}

	// max page:
	pager.MaxPage = count/pager.PageSize + 1

	// page:
	pager.Current = currentPage
	if currentPage < 1 {
		pager.Current = 1
	} else if currentPage > pager.MaxPage {
		pager.Current = pager.MaxPage
	}

	return pager
}
