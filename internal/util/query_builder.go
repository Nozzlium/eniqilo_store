package util

import (
	"bytes"
	"fmt"
)

func BuildQueryStringAndParams(
	baseQuery *bytes.Buffer,
	whereBuilder func() ([]string, []interface{}),
	paginationBuilder func() (string, []interface{}),
	orderByBuilder func() ([]string, []interface{}),
) (string, []interface{}) {
	where, params := whereBuilder()
	for i, clause := range where {
		fmt.Fprintf(baseQuery, " and %s", fmt.Sprintf(clause, i+1))
	}

	pagination, paginationParams := paginationBuilder()
	fmt.Fprintf(baseQuery, fmt.Sprintf(" %s", pagination), len(params)+1, len(params)+2)
	params = append(params, paginationParams...)

	orderBy, orderByParams := orderByBuilder()
	if len(orderByParams) > 0 {
		baseQuery.WriteString(" order by")
		params = append(params, orderByParams...)
		for i, clause := range orderBy {
			if i > 1 {
				baseQuery.WriteString(",")
			}
			fmt.Fprintf(baseQuery, " %s", fmt.Sprintf(clause, i))
		}
	}

	return baseQuery.String(), params
}
