package mysql_test

import (
	"context"
	"testing"
	"time"

	mysql_ "github.com/kaydxh/golang/pkg/database/mysql"
	viper_ "github.com/kaydxh/golang/pkg/viper"
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
			Address:  "127.0.0.1:3306",
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

func TestGetDatabaseUntil(t *testing.T) {
	testCases := []struct {
		Address  string
		DataName string
		UserName string
		Password string
	}{
		{
			Address:  "127.0.0.1:3306",
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
			sqlDB, err := db.GetDatabaseUntil(context.Background(), 5*time.Second, 20*time.Second)
			if err != nil {
				t.Fatalf("failed to get database: %v, got : %s", testCase.DataName, err)
			}
			assert.NotNil(t, sqlDB)
		})
	}

}

func TestNew(t *testing.T) {

	cfgFile := "./mysql.yaml"
	config := mysql_.NewConfig(mysql_.WithViper(viper_.GetViper(cfgFile, "database.mysql")))

	db, err := config.Complete().New(context.Background())
	if err != nil {
		t.Errorf("failed to new config err: %v", err)
	}

	t.Logf("db: %#v", db)
}

func TestGetTheDBAndClose(t *testing.T) {
	testCases := []struct {
		Address  string
		DataName string
		UserName string
		Password string
	}{
		{
			Address:  "127.0.0.1:3306",
			DataName: "test",
			UserName: "root",
			Password: "123456",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Address+testCase.DataName, func(t *testing.T) {
			conf := mysql_.DBConfig{
				Address:  testCase.Address,
				DataName: testCase.DataName,
				UserName: testCase.UserName,
				Password: testCase.Password,
			}
			db := mysql_.NewDB(conf)
			sqlDB, err := db.GetDatabaseUntil(context.Background(), 5*time.Second, 20*time.Second)
			if err != nil {
				t.Fatalf("failed to get database: %v, got : %s", testCase.DataName, err)
			}
			assert.NotNil(t, sqlDB)

			theDB, err := mysql_.GetTheDB(conf)
			assert.Nil(t, err)

			assert.Equal(t, sqlDB, mysql_.GetDB())
			t.Logf("GetDB got sqlDB: %v, expect %v", sqlDB, mysql_.GetDB())
			t.Logf("GetTheDB got sqlDB: %v, expect %v", sqlDB, theDB)
			err = mysql_.CloseTheDB(conf)
			assert.Nil(t, err)
			err = mysql_.CloseDB()
			assert.Nil(t, err)
		})
	}

}
