package mysql

import (
	"context"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	reflect_ "github.com/kaydxh/golang/go/reflect"
)

func GetByQuery(ctx context.Context, db *sqlx.DB, query string, arg interface{}) (interface{}, error) {
	// Check that invalid preparations fail
	ns, err := db.PrepareNamedContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer ns.Close()

	var dest interface{}
	err = ns.SelectContext(ctx, &dest, arg)
	if err != nil {
		return nil, err
	}
	return dest, nil
}

/*
type SqlCompare int

const (
	SqlCompareEqual      SqlCompare = iota //=
	SqlCompareNotEqual   SqlCompare = iota //<>
	SqlCompareGreater    SqlCompare = iota //>
	SqlCompareLessThan   SqlCompare = iota //<
	SqlCompareGreatEqual SqlCompare = iota //>=
	SqlCompareLessEqual  SqlCompare = iota //<=
	SqlCompareLike       SqlCompare = iota //LIKE
)
*/
type SqlCompare string

const (
	SqlCompareEqual      SqlCompare = "="
	SqlCompareNotEqual   SqlCompare = "!="
	SqlCompareGreater    SqlCompare = ">"
	SqlCompareLessThan   SqlCompare = "<"
	SqlCompareGreatEqual SqlCompare = ">="
	SqlCompareLessEqual  SqlCompare = "<="
	SqlCompareLike       SqlCompare = "LIKE"
)

type SqlOperator string

const (
	SqlOperatorAnd SqlOperator = "AND"
	SqlOperatorOr  SqlOperator = "OR"
	SqlOperatorNot SqlOperator = "NOT"
)

func GenerateCondition(cmp SqlCompare, query string, arg interface{}) string {
	condFields := reflect_.NonzeroFields(arg)
	return JoinNamedTableColumnsValues(cmp, condFields...)
}

//  "foo=:foo, bar=:bar"
func JoinNamedTableColumnsValues(cmp SqlCompare, cols ...string) string {
	return strings.Join(NamedTableColumnsValues(SqlCompareEqual, cols...), ",")
}

// NamedColumnsValues returns the []string{table.value1=:value1, table.value2=:value2 ...}
// query := NamedColumnsValues("table", "foo", "bar")
// // []string{"table.foo=:table.foo", "table.bar=:table.bar"}
func NamedTableColumnsValues(cmp SqlCompare, cols ...string) []string {

	var namedCols []string
	for _, col := range cols {
		namedCols = append(namedCols, fmt.Sprintf("%[1]s %[2]s :%[1]s", col, cmp))
	}
	return namedCols
}
