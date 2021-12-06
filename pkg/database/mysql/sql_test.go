package mysql_test

import (
	"fmt"
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

func TestNamedInCondition(t *testing.T) {

	testCases := []struct {
		cols []string
		arg  interface{}
	}{
		{
			cols: []string{"task_id"},
			arg: struct {
				TaskId []string `db:"task_id"`
			}{
				TaskId: []string{"task_id_1", "task_id_2"},
			},
		},
		{
			cols: []string{"task_id"},
			arg: struct {
				TaskId []int `db:"task_id"`
			}{
				TaskId: []int{0, 1},
			},
		},
	}

	for i, testCase := range testCases {
		t.Run(fmt.Sprintf("%v", i), func(t *testing.T) {
			query, err := mysql_.NamedInCondition(mysql_.SqlOperatorAnd, testCase.cols, testCase.arg)
			if err != nil {
				t.Fatalf("err: %v", err)
			}
			t.Logf("sql: %v", query)
		})
	}
}

/*
func TestGenerateSQL(t *testing.T) {
	arg := struct {
		TaskId string   `db:"task_id"`
		Name   []string `db:"name"`
	}{
		TaskId: "task-id",
		Name:   []string{"name1", "name2"},
	}
	//sql := "SELECT *FROM t_task where task_id=:task_id and name In(:name)"
	sql := "SELECT *FROM t_task where task_id=:task_id"
	query, args, err := sqlx.Named(sql, arg)
	if err != nil {
		t.Fatalf("falied to named, err: %v", err)
	}
	query, args, err = sqlx.In(query, args...)
	if err != nil {
		t.Fatalf("falied to In, err: %v", err)
	}

	// ns, err := d.db.PrepareNamedContext(ctx, query)
	t.Logf("query: %v, args: %v", query, args)
}
*/
