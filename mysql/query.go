package mysql

import (
	"fmt"
)

type Query struct {
	Sql        string `db:"-" json:"-"`
	Where      string `db:"-" json:"-"`
	AfterWhere string `db:"-" json:"-"`
}

func (q *Query) Form(tableName string) *Query {
	q.Sql = fmt.Sprintf("%s FROM `%s`", q.Sql, tableName)
	return q
}

func (q *Query) LeftJoin(tableName, on string) *Query {
	q.Sql = fmt.Sprintf("%s LEFT JOIN (`%s`) ON (%s)", q.Sql, tableName, on)
	return q
}

func (q *Query) OrderBy(field string) *Query {
	q.AfterWhere = fmt.Sprintf("%s ORDER BY `%s`", q.AfterWhere, field)
	return q
}

func (q *Query) OrderAsc(field string) *Query {
	q.AfterWhere = fmt.Sprintf("%s ORDER BY `%s` ASC", q.AfterWhere, field)
	return q
}

func (q *Query) OrderDesc(field string) *Query {
	q.AfterWhere = fmt.Sprintf("%s ORDER BY `%s` DESC", q.AfterWhere, field)
	return q
}

func (q *Query) Limit(limit uint64) *Query {
	q.AfterWhere = fmt.Sprintf("%s LIMIT %d", q.AfterWhere, limit)
	return q
}

func (q *Query) LimitPage(offset, limit uint64) *Query {
	q.AfterWhere = fmt.Sprintf("%s LIMIT %d,%d", q.AfterWhere, offset, limit)
	return q
}

func (q *Query) Combination() string {
	cmd := q.Sql
	if q.Where != "" {
		cmd = cmd + q.Where

	}
	if q.AfterWhere != "" {
		cmd = cmd + q.AfterWhere
	}
	return cmd
}
