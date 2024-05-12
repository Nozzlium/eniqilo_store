package util

import (
	"bytes"
	"fmt"
)

func BuildQueryStringAndParams(
	baseQuery *bytes.Buffer,
	whereBuilder func() ([]string, []interface{}),
	paginationBuilder func() (string, []interface{}),
	orderByBuilder func() []string,
) (string, []interface{}) {
	where, params := whereBuilder()
	for i, clause := range where {
		fmt.Fprintf(baseQuery, " and %s", fmt.Sprintf(clause, i+1))
	}

	orderBy := orderByBuilder()
	if len(orderBy) > 0 {
		baseQuery.WriteString(" order by ")
		for i, clause := range orderBy {
			if i > 0 {
				baseQuery.WriteString(", ")
			}
			fmt.Fprint(baseQuery, clause)
		}
	}

	pagination, paginationParams := paginationBuilder()
	paramsLen := len(params)
	fmt.Fprintf(baseQuery, fmt.Sprintf(" %s ", pagination), paramsLen+1, paramsLen+2)
	params = append(params, paginationParams...)

	return baseQuery.String(), params
}

func BuildQueryStringAndParamsWithoutLimit(
	baseQuery *bytes.Buffer,
	whereBuilder func() ([]string, []interface{}),
	orderByBuilder func() []string,
) (string, []interface{}) {
	where, params := whereBuilder()
	for i, clause := range where {
		fmt.Fprintf(baseQuery, " and %s", fmt.Sprintf(clause, i+1))
	}

	orderBy := orderByBuilder()
	if len(orderBy) > 0 {
		baseQuery.WriteString(" order by ")
		for i, clause := range orderBy {
			if i > 0 {
				baseQuery.WriteString(", ")
			}
			fmt.Fprint(baseQuery, clause)
		}
	}

	return baseQuery.String(), params
}

func DefaultPaginationBuilder(limit, offset int) (string, []interface{}) {

	defaultLimit := 5
	defaultOffset := 0
	if limit > 0 {
		defaultLimit = limit
	}
	if offset > 0 {
		defaultOffset = offset
	}

	return " limit $%d offset $%d ", []interface{}{defaultLimit, defaultOffset}
}
