package mysql

import (
	"fmt"
	"strings"

	reflect_ "github.com/kaydxh/golang/go/reflect"
)

const dbTag = "db"

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

// "WHERE foo=:foo AND bar=:bar"
func GenerateCondition(cmp SqlCompare, oper SqlOperator, query string, arg interface{}) string {
	condFields := reflect_.NonzeroFieldTags(arg, dbTag)
	return fmt.Sprintf("%s %s", query, func() string {
		if len(condFields) == 0 {
			return ""
		}
		return fmt.Sprintf(" WHERE %s", JoinNamedColumnsValuesWithOperator(cmp, oper, condFields...))
	}())
}

// "WHERE foo=:foo AND bar=:bar"
func GenerateNameColumsCondition(cmp SqlCompare, oper SqlOperator, condFields ...string) string {
	return fmt.Sprintf(" %s ", func() string {
		if len(condFields) == 0 {
			return ""
		}
		return fmt.Sprintf(" WHERE %s", JoinNamedColumnsValuesWithOperator(cmp, oper, condFields...))
	}())
}

//foo=:foo,bar=:bar,  for update set
func JoinNamedColumnsValues(cols ...string) string {
	return strings.Join(namedTableColumnsValues(SqlCompareEqual, cols...), ",")
}

// "foo=:foo AND bar=:bar" , for where condition
func JoinNamedColumnsValuesWithOperator(cmp SqlCompare, oper SqlOperator, cols ...string) string {
	return strings.Join(namedTableColumnsValues(cmp, cols...), fmt.Sprintf(" %s ", oper))
}

// []string{"foo=:foo",  "bar=:bar"}
func namedTableColumnsValues(cmp SqlCompare, cols ...string) []string {
	var namedCols []string
	for _, col := range cols {
		namedCols = append(namedCols, fmt.Sprintf("%[1]s %[2]s :%[1]s", col, cmp))
	}
	return namedCols
}