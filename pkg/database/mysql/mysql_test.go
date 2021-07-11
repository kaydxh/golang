package mysql_test

import (
	"testing"

	mysql_ "github.com/kaydxh/golang/pkg/database/mysql"
	"github.com/stretchr/testify/assert"
)

func TestGetDataBase(t *testing.T) {
	testCases := []struct {
		Address  string
		DataName string
		UserName string
		Password string
	}{
		{
			Address:  "127.0.0.1",
			DataName: "test",
			UserName: "root",
			Password: "123456",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Address+testCase.DataName, func(t *testing.T) {
			db := mysql_.NewDB(mysql_.DBConfig{
				Address:  testCase.Address,
				DataName: testCase.DataName,
				UserName: testCase.UserName,
				Password: testCase.Password,
			})
			sqlDB, err := db.GetDatabase()
			if err != nil {
				t.Fatalf("failed to get database: %v, got : %s", testCase.DataName, err)
			}
			assert.NotNil(t, sqlDB)
		})
	}

}
