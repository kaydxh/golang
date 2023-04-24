/*
 *Copyright (c) 2023, kaydxh
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
package mysql_test

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"testing"

	"github.com/google/uuid"
	mysql_ "github.com/kaydxh/golang/pkg/database/mysql"
	viper_ "github.com/kaydxh/golang/pkg/viper"
)

type TaskTable struct {
	//	Id sql.NullInt64 `db:"id"` // primary key ID

	// NullTime represents a time.Time that may be null.
	// NullTime implements the Scanner interface so
	// it can be used as a scan destination, similar to NullString.
	CreateTime sql.NullTime `db:"create_time"`
	UpdateTime sql.NullTime `db:"update_time"`

	GroupId    string `db:"group_id"`
	PageId     string `db:"page_id"`
	FeaId      string `db:"fea_id"`
	EntityId   string `db:"entity_id"`
	Feature0   []byte `db:"feature0"`
	Feature1   []byte `db:"feature1"`
	ExtendInfo []byte `db:"extend_info"`
}

type Tasks []*TaskTable

func (t Tasks) String() string {
	s := "["
	for _, task := range t {
		s += fmt.Sprintf("%v,", task)
	}
	if len(t) > 0 {
		s = strings.TrimRight(s, ",")
	}

	return s + "]"
}
func TestInsert(t *testing.T) {

	cfgFile := "./mysql.dev.yaml"
	config := mysql_.NewConfig(mysql_.WithViper(viper_.GetViper(cfgFile, "database.mysql")))

	db, err := config.Complete().New(context.Background())
	if err != nil {
		t.Errorf("failed to new config err: %v", err)
	}

	t.Logf("db: %#v", db)
	ctx := context.Background()

	testCases := []struct {
		TableName string
		GroupId   string
		Number    int64
	}{
		{
			TableName: "hetu_zeus_0",
			GroupId:   "hetu_image_copyright_prod",
			Number:    100,
		},
	}

	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("table[%v]-fieldKey[%v]", testCase.TableName, testCase.GroupId), func(t *testing.T) {

			count := testCase.Number

			for count > 0 {
				arg := &TaskTable{
					GroupId:    testCase.GroupId,
					PageId:     "100",
					FeaId:      uuid.NewString(),
					Feature0:   []byte("Feature0"),
					Feature1:   []byte("Feature1"),
					ExtendInfo: []byte("ExtendInfo"),
				}
				query := fmt.Sprintf(`INSERT INTO %s
			   (
			   group_id,
			   page_id,
			   fea_id,
			   entity_id,
			   feature0,
			   feature1,
			   extend_info
			   )
			   VALUES (
			         :group_id,
			         :page_id,
					 :fea_id,
					 :entity_id,
					 :feature0,
					 :feature1,
					 :extend_info
					 )
                     `, testCase.TableName)

				_, err = mysql_.ExecContext(ctx, query, arg, nil, db)
				if err != nil {
					t.Errorf("faild to insert %v, err: %v", arg, err)
				}
				count--
			}

		})
	}
}

func TestInsertNewBatch(t *testing.T) {

	cfgFile := "./mysql.dev.yaml"
	config := mysql_.NewConfig(mysql_.WithViper(viper_.GetViper(cfgFile, "database.mysql")))

	db, err := config.Complete().New(context.Background())
	if err != nil {
		t.Errorf("failed to new config err: %v", err)
	}

	t.Logf("db: %#v", db)
	ctx := context.Background()

	batch := 512
	tableName := "hetu_zeus_0"

	cols := []string{"group_id", "page_id", "fea_id", "entity_id", "feature0", "feature1", "extend_info"}
	query := fmt.Sprintf(`INSERT INTO %s
			   (
			   group_id,
			   page_id,
			   fea_id,
			   entity_id,
			   feature0,
			   feature1,
			   extend_info
			   )
			   VALUES %s
                     `, tableName,
		mysql_.JoinNamedColumnsValuesBatch(cols, batch))

	testCases := []struct {
		TableName string
		GroupId   string
		Number    int64
	}{
		{
			TableName: "hetu_zeus_0",
			GroupId:   "hetu_image_copyright_prod",
		},
	}

	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("table[%v]-fieldKey[%v]", testCase.TableName, testCase.GroupId), func(t *testing.T) {

			count := batch
			var args []*TaskTable
			for count > 0 {
				arg := &TaskTable{
					GroupId:    testCase.GroupId,
					PageId:     "100",
					FeaId:      uuid.NewString(),
					EntityId:   "200",
					Feature0:   []byte("Feature0"),
					Feature1:   []byte("Feature1"),
					ExtendInfo: []byte("ExtendInfo"),
				}
				args = append(args, arg)
				count--
			}
			_, err = mysql_.ExecContext(ctx, query, mysql_.BuildNamedColumnsValuesBatch(args), nil, db)
			if err != nil {
				t.Errorf("faild to insert %v, err: %v", args, err)
			}

		})
	}
}

func TestInsertBatch(t *testing.T) {

	cfgFile := "./mysql.dev.yaml"
	config := mysql_.NewConfig(mysql_.WithViper(viper_.GetViper(cfgFile, "database.mysql")))

	db, err := config.Complete().New(context.Background())
	if err != nil {
		t.Errorf("failed to new config err: %v", err)
	}

	t.Logf("db: %#v", db)
	ctx := context.Background()

	tableName := "hetu_zeus_0"
	query := fmt.Sprintf(`INSERT INTO %s
			   (
			   group_id,
			   page_id,
			   fea_id,
			   entity_id,
			   feature0,
			   feature1,
			   extend_info
			   )
			   VALUES (
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
                     `, tableName)

	testCases := []struct {
		TableName string
		GroupId   string
		Number    int64
	}{
		{
			TableName: "hetu_zeus_0",
			GroupId:   "hetu_image_copyright_prod",
			Number:    2,
		},
	}

	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("table[%v]-fieldKey[%v]", testCase.TableName, testCase.GroupId), func(t *testing.T) {

			count := testCase.Number
			var args []TaskTable

			for count > 0 {
				arg := TaskTable{
					GroupId:    testCase.GroupId,
					PageId:     "100",
					FeaId:      uuid.NewString(),
					EntityId:   "200",
					Feature0:   []byte("Feature0"),
					Feature1:   []byte("Feature1"),
					ExtendInfo: []byte("ExtendInfo"),
				}
				args = append(args, arg)
				count--
			}

			tagsValues := map[string]interface{}{
				"group_id_1":    "group_id_1",
				"page_id_1":     uuid.NewString(),
				"fea_id_1":      uuid.NewString(),
				"entity_id_1":   uuid.NewString(),
				"feature0_1":    "feature0_1",
				"feature1_1":    "feature1_1",
				"extend_info_1": "extend_info_1",
				"group_id_2":    "group_id_2",
				"page_id_2":     uuid.NewString(),
				"fea_id_2":      uuid.NewString(),
				"entity_id_2":   uuid.NewString(),
				"feature0_2":    "feature0_2",
				"feature1_2":    "feature1_2",
				"extend_info_2": "extend_info_2",
			}
			t.Logf("batch insert %v", tagsValues)

			//_, err = mysql_.ExecContext(ctx, query, tagsValues, nil, db)
			//_, err = mysql_.ExecContext(ctx, query, []interface{}{tagsValues[0], tagsValues[1]}, nil, db)
			_, err = mysql_.ExecContext(ctx, query, tagsValues, nil, db)
			if err != nil {
				t.Errorf("faild to insert %v, err: %v", tagsValues, err)
			}

		})
	}
}

func TestDeleteBatch(t *testing.T) {

	cfgFile := "./mysql.yaml"
	config := mysql_.NewConfig(mysql_.WithViper(viper_.GetViper(cfgFile, "database.mysql")))

	db, err := config.Complete().New(context.Background())
	if err != nil {
		t.Errorf("failed to new config err: %v", err)
	}

	t.Logf("db: %#v", db)
	ctx := context.Background()

	testCases := []struct {
		TableName   string
		GroupId     string
		DeleteField string
		Batch       int64
	}{
		{
			TableName:   "hetu_zeus_0",
			DeleteField: "group_id",
			GroupId:     "hetu_image_copyright_prod",
			Batch:       5,
		},
	}

	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("table[%v]-fieldKey[%v]", testCase.TableName, testCase.GroupId), func(t *testing.T) {

			var count int64
			arg := &TaskTable{
				GroupId: testCase.GroupId,
			}

			for {
				query := fmt.Sprintf(
					`DELETE FROM %s
	           WHERE %s limit %v`,
					testCase.TableName,
					mysql_.ConditionWithEqualAnd(testCase.DeleteField),
					testCase.Batch,
				)
				rows, err := mysql_.ExecContext(ctx, query, arg, nil, db)
				if err != nil {
					t.Fatalf("failed to delete %v, current deleted total number: %v, err: %v", arg.GroupId, count, err)
				}

				count += rows
				if rows == 0 {
					t.Logf("finished to delete %v, total number: %v", arg.GroupId, count)
					break
				}

				if count%testCase.Batch == 0 {
					t.Logf("delete number: %v ...", count)
				}
			}

		})
	}

}
