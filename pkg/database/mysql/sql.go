package mysql

import (
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	reflect_ "github.com/kaydxh/golang/go/reflect"
	strings_ "github.com/kaydxh/golang/go/strings"
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
	SqlCompareIn         SqlCompare = "IN"
)

type SqlOperator string

const (
	SqlOperatorAnd SqlOperator = "AND"
	SqlOperatorOr  SqlOperator = "OR"
	SqlOperatorNot SqlOperator = "NOT"
)

// "foo=:foo AND bar=:bar"
func NonzeroCondition(cmp SqlCompare, oper SqlOperator, arg interface{}) string {
	condFields := reflect_.NonzeroFieldTags(arg, dbTag)
	return fmt.Sprintf(" %s ", func() string {
		if len(condFields) == 0 {
			return "TRUE"
		}
		return fmt.Sprintf("%s", JoinNamedColumnsValuesWithOperator(cmp, oper, condFields...))
	}())
}

func NonzeroFields(arg interface{}) []string {
	return reflect_.NonzeroFieldTags(arg, dbTag)
}

func ConditionWithEqualAnd(condFields ...string) string {
	return JoinNamedColumnsValuesWithOperator(SqlCompareEqual, SqlOperatorAnd, condFields...)
}

// "ORDER BY create_time DESC, id DESC"
func OrderCondition(orders map[string]bool) string {
	if len(orders) == 0 {
		return ""
	}

	return fmt.Sprintf(" ORDER BY %s", func() string {
		var msg string
		for k, v := range orders {
			msg += fmt.Sprintf("%s %s,", k, func() string {
				if v {
					return "DESC"
				}
				return "ASC"
			}())
		}

		msg = strings.TrimRight(msg, ",")
		return msg
	}())

}

func InCondition(cond string, values ...string) string {
	if cond == "" || len(values) == 0 {
		return "TRUE"
	}

	return fmt.Sprintf(`%s IN (%s)`, cond, func() string {
		var msg string
		for _, v := range values {
			msg += fmt.Sprintf(`"%s",`, v)
		}
		msg = strings.TrimRight(msg, ",")
		return msg
	}())
}

func NamedInCondition(oper SqlOperator, cols []string, arg interface{}) (string, error) {
	query := JoinNamedColumnsValuesWithOperator(SqlCompareIn, oper, cols...)
	query, args, err := sqlx.Named(query, arg)
	if err != nil {
		return "", err
	}
	query, args, err = sqlx.In(query, args...)
	if err != nil {
		return "", err
	}

	return strings_.ReplaceAll(query, "?", args, true), nil
}

//foo=:foo,bar=:bar,  for update set
func JoinNamedColumnsValues(cols ...string) string {
	return strings.Join(namedTableColumnsValues(SqlCompareEqual, cols...), ",")
}

// "foo=:foo AND bar=:bar" , for where condition
func JoinNamedColumnsValuesWithOperator(cmp SqlCompare, oper SqlOperator, cols ...string) string {
	conds := strings.Join(namedTableColumnsValues(cmp, cols...), fmt.Sprintf(" %s ", oper))
	if len(cols) == 0 || conds == "" {
		return "TRUE"
	}

	return conds
}

// []string{"foo=:foo",  "bar=:bar"}
func namedTableColumnsValues(cmp SqlCompare, cols ...string) []string {
	var namedCols []string
	for _, col := range cols {
		if col != "" {
			switch cmp {
			case SqlCompareLike:
				namedCols = append(namedCols, fmt.Sprintf(`%[1]s %[2]s concat("%%",:%[1]s,"%%")`, col, cmp))
			case SqlCompareIn:
				namedCols = append(namedCols, fmt.Sprintf("%[1]s %[2]s (:%[1]s)", col, cmp))
			default:
				namedCols = append(namedCols, fmt.Sprintf("%[1]s %[2]s :%[1]s", col, cmp))
			}
		}
	}
	return namedCols
}
