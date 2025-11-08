package service

type MandatoryRequest struct {
	AcceptLanguage string
	Authorization  string
}

type BaseEntity struct {
	Id int64 `db:"ID" json:"id" example:"12345"`
}

type SqlParameter struct {
	TableName string `json:"table,omitempty"`

	// selected columns for select query.
	Columns []string `json:"col,omitempty"`

	// field, value pair for insert / update.
	Values []Value `json:"values,omitempty"`

	// filter for WHERE clause
	Params []FilterParam `json:"params,omitempty"`

	// Order By Clause (column ASC / DESC)
	OrderBy []string `json:"order,omitempty"`

	// Having
	Having []string `json:"having,omitempty"`

	// Limit and Offset for pagination
	Limit  int `json:"limit,omitempty"`
	Offset int `json:"offset,omitempty"`

	// Export parameters
	GroupBy    []string `json:"group,omitempty"`
	ExportType string   `json:"export,omitempty"`

	// extra parameters for passing specific params to repository
	CacheKey string `json:"cacheKey,omitempty"`

	// Extra parameters
	Extra map[string]string `json:"extra,omitempty"`

	// New: List of joins
	Joins []JoinClause `json:"joins,omitempty"`

	SelectField string `json:"-"` //optional,if want to put here new field first need adding comma  e.g. ,name
}

type Value struct {
	Field string      `json:"field,omitempty"`
	Value interface{} `json:"value,omitempty"`
}

type PaginationParam struct {
	TableName       string
	Columns         []string
	Limit           int
	Last            string
	ExtraParamKey   []string
	ExtraParamValue []interface{}
	Offset          int
}

type FilterParam struct {
	Field   string      `json:"field"`
	Operand string      `json:"operand"`
	Value   interface{} `json:"value"`
}

type JoinClause struct {
	Table      string        `json:"table"`
	Alias      string        `json:"alias,omitempty"`
	On         string        `json:"on"`                   // raw ON condition
	JoinType   string        `json:"joinType"`             // e.g. "LEFT", "INNER"
	Conditions []FilterParam `json:"conditions,omitempty"` // optional
}
type Pagination struct {
	LastDate  string `json:"last_date,omitempty" example:"2021-05-31T08:20:02Z"`
	TotalData int64  `json:"totalData,omitempty"`
	NextPage  bool   `json:"nextPage"`
	LastID    int64  `json:"lastId,omitempty"`
}

func MakePagination(count int64, params SqlParameter, lenData int) Pagination {
	next := true
	if lenData == 0 {
		next = false
	}
	if int64(params.Offset+lenData) >= count {
		next = false
	}
	return Pagination{
		TotalData: count,
		NextPage:  next,
	}
}
