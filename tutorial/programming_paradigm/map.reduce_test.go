package tutorial

import (
	"fmt"
	"testing"
)

type Employee struct {
	Name     string
	Age      int
	Vacation int
	Salary   int
}

var list = []Employee{
	{"Hao", 44, 0, 8000},
	{"Bob", 34, 10, 5000},
	{"Alice", 23, 5, 9000},
	{"Jack", 26, 0, 4000},
	{"Tom", 48, 9, 7500},
	{"Marry", 29, 0, 6000},
	{"Mike", 32, 8, 4000},
}

// 控制逻辑通过传递函数方式与业务逻辑分开
func EmployeeCountIf(list []Employee, fn func(e *Employee) bool) int {
	count := 0
	for i, _ := range list {
		if fn(&list[i]) {
			count += 1
		}
	}
	return count
}

func EmployeeFilterIn(list []Employee, fn func(e *Employee) bool) []Employee {
	var newList []Employee
	for i, _ := range list {
		if fn(&list[i]) {
			newList = append(newList, list[i])
		}
	}
	return newList
}

func EmployeeSumIf(list []Employee, fn func(e *Employee) int) int {
	var sum = 0
	for i, _ := range list {
		sum += fn(&list[i])
	}
	return sum
}

func TestEmployeeCountIf(t *testing.T) {
	testCases := []struct {
		Age    int
		Salary int
	}{
		{
			Age:    40,
			Salary: 6000,
		},
	}

	for i, testCase := range testCases {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			count := EmployeeCountIf(list, func(e *Employee) bool {
				return e.Age > testCase.Age
			})
			t.Logf("the number of age > %v is %v", testCase.Age, count)
			count = EmployeeCountIf(list, func(e *Employee) bool {
				return e.Salary > testCase.Salary
			})
			t.Logf("the number of salary > %v is %v", testCase.Salary, count)
		})
	}
}
