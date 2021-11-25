package mysql_test

import (
	"testing"

	mysql_ "github.com/kaydxh/golang/pkg/database/mysql"
)

func TestJoinNamedColumnsValuesWithOperator(t *testing.T) {
	testCases := []struct {
		cmp  mysql_.SqlCompare
		oper mysql_.SqlOperator
		cols []string
	}{
		{
			cmp:  mysql_.SqlCompareLike,
			oper: mysql_.SqlOperatorAnd,
			cols: []string{"task_name"},
		},
	}

	for _, testCase := range testCases {
		t.Run(string(testCase.cmp), func(t *testing.T) {
			query := mysql_.JoinNamedColumnsValuesWithOperator(testCase.cmp, testCase.oper, testCase.cols...)
			t.Logf("sql: %v", query)
		})
	}
}

func TestGenerateInCondition(t *testing.T) {
	testCases := []struct {
		cond   string
		values []string
	}{
		{
			cond:   "task_id",
			values: []string{"task_id_1", "task_id_2"},
		},
		{
			cond:   "task_id",
			values: []string{"", ""},
		},
	}

	for _, testCase := range testCases {
		t.Run(string(testCase.cond), func(t *testing.T) {
			query := mysql_.GenerateInCondition(testCase.cond, testCase.values...)
			t.Logf("sql: %v", query)
		})
	}
}
