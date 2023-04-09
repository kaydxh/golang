/*
 *Copyright (c) 2022, kaydxh
 *
 *Permission is hereby granted, free of charge, to any person obtaining a copy
 *of this software and associated documentation files (the "Software"), to deal
 *in the Software without restriction, including without limitation the rights
 *to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 *copies of the Software, and to permit persons to whom the Software is
 *furnished to do so, subject to the following conditions:
 *
 *The above copyright notice and this permission notice shall be included in all
 *copies or substantial portions of the Software.
 *
 *THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 *IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 *FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 *AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 *LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 *OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 *SOFTWARE.
 */
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

// JoinNamedColumnsValues foo=:foo,bar=:bar,  for update set
func JoinNamedColumnsValues(cols ...string) string {
	return strings.Join(namedTableColumnsValues(SqlCompareEqual, cols...), ",")
}

// JoinNamedColumnsValuesWithOperator "foo=:foo AND bar=:bar" , for where condition
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

/*
used for batch insert
(
 :group_id_1,
 :page_id_1,
 :fea_id_1,
 :entity_id_1,
 :feature0_1,
 :feature1_1,
 :extend_info_1
 ),
(
 :group_id_2,
 :page_id_2,
 :fea_id_2,
 :entity_id_2,
 :feature0_2,
 :feature1_2,
 :extend_info_2
 )
*/
func JoinNamedColumnsValuesBatch(cols []string, batch int) string {

	var batchNamedCols []string
	for i := 0; i < batch; i++ {
		var namedCols []string
		for _, col := range cols {
			namedCols = append(namedCols, fmt.Sprintf(":%s_%d", col, i))
		}
		batchNamedCols = append(batchNamedCols, fmt.Sprintf("(%v)", strings.Join(namedCols, ",")))
	}

	return strings.Join(batchNamedCols, ",")
}

// used for batch insert
func TransferToNamedColumnsValuesBatch(req []map[string]interface{}) map[string]interface{} {

	valuesMap := make(map[string]interface{}, 0)
	for i, values := range req {
		for k, v := range values {
			valuesMap[fmt.Sprintf("%s_%d", k, i)] = v
		}
	}

	return valuesMap
}

// req is slice of struct or pointer struct
func BuildNamedColumnsValuesBatch(req interface{}) map[string]interface{} {
	return TransferToNamedColumnsValuesBatch(reflect_.ArrayAllTagsVaules(req, dbTag))
}
