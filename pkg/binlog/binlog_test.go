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

package binlog_test

import (
	"context"
	"database/sql"
	"testing"

	"encoding/json"

	"github.com/google/uuid"
	"github.com/kaydxh/golang/pkg/binlog"
	binlog_ "github.com/kaydxh/golang/pkg/binlog"
	ds_ "github.com/kaydxh/golang/pkg/binlog/datastore"
	mysql_ "github.com/kaydxh/golang/pkg/database/mysql"
	mq_ "github.com/kaydxh/golang/pkg/mq"
	kafka_ "github.com/kaydxh/golang/pkg/mq/kafka"
	viper_ "github.com/kaydxh/golang/pkg/viper"
	"github.com/segmentio/kafka-go"
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

func TestProducer(t *testing.T) {
	cfgFile := "./binlog.yaml"
	config := kafka_.NewConfig(kafka_.WithViper(viper_.GetViper(cfgFile, "mq.kafka")))

	mq, err := config.Complete().New(context.Background())
	if err != nil {
		t.Errorf("failed to new config err: %v", err)
	}

	topic := "topic-test-1"
	ctx := context.Background()
	err = mq.AsProducers(ctx, topic)
	if err != nil {
		t.Fatalf("failed to as producers, err: %v", err)
	}
	p, err := mq.GetProducer(topic)
	if err != nil {
		t.Fatalf("failed to get producer, err: %v", err)
	}

	tableName := "hetu_zeus_0"
	cols := []string{"group_id", "page_id", "fea_id", "entity_id", "feature0", "feature1", "extend_info"}
	msgKey := &ds_.MessageKey{
		Key:     "Key-A",
		MsgType: ds_.MsgType_Insert,
		Fields:  cols,
		Path:    tableName,
	}
	keyData, _ := json.Marshal(msgKey)

	for i := 0; i < 10; i++ {
		arg := &TaskTable{
			GroupId:    "groupId-1",
			PageId:     "100",
			FeaId:      uuid.NewString(),
			Feature0:   []byte("Feature0"),
			Feature1:   []byte("Feature1"),
			ExtendInfo: []byte("ExtendInfo"),
		}
		msgValue, _ := json.Marshal(arg)
		err = p.Send(ctx,
			kafka.Message{
				Key:   keyData,
				Value: msgValue,
			},
		)
		if err != nil {
			t.Fatalf("failed to send messages, err: %v", err)
		}
	}

}

func TestNewBinlog(t *testing.T) {
	// install kafka
	cfgFile := "./binlog.yaml"
	config := kafka_.NewConfig(kafka_.WithViper(viper_.GetViper(cfgFile, "mq.kafka")))

	mq, err := config.Complete().New(context.Background())
	if err != nil {
		t.Errorf("failed to new config err: %v", err)
	}

	ctx := context.Background()
	topic := "topic-test-1"
	err = mq.AsConsumers(ctx, topic)
	if err != nil {
		t.Fatalf("failed to as producers, err: %v", err)
	}
	consumer, err := mq.GetConsumer(topic)
	if err != nil {
		t.Fatalf("failed to get producer, err: %v", err)
	}

	dbConfig := mysql_.NewConfig(mysql_.WithViper(viper_.GetViper(cfgFile, "database.mysql")))

	db, err := dbConfig.Complete().New(context.Background())
	if err != nil {
		panic(err)
	}
	if db == nil {
		panic("db is not enable")
	}

	bsConfig := binlog_.NewConfig(binlog_.WithViper(viper_.GetViper(cfgFile, "binlog")))

	bs, err := bsConfig.Complete().New(context.Background(), nil, []mq_.Consumer{consumer},
		binlog.WithMessageDecoderFunc(func(ctx context.Context, data []byte) (interface{}, error) {

			var arg TaskTable
			err = json.Unmarshal(data, &arg)
			if err != nil {
				return arg, err
			}
			return arg, err

		}),

		binlog_.WithMessageKeyDecodeFunc(func(ctx context.Context, data []byte) (ds_.MessageKey, error) {
			var msgKey ds_.MessageKey
			err := json.Unmarshal(data, &msgKey)
			return msgKey, err

		}))
	if err != nil {
		panic(err)
	}
	_ = bs

	select {}
}
