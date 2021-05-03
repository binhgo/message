package model

import (
	"time"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"gitlab.ghn.vn/common-projects/go-sdk/sdk"
	"gitlab.ghn.vn/internal-tools/message/model/enum"
)

type MessageLog struct {
	ID              bson.ObjectId `json:"id,omitempty" bson:"_id,omitempty"`
	LastUpdatedTime *time.Time    `json:"lastUpdatedTime,omitempty" bson:"last_updated_time,omitempty"`
	CreatedTime     *time.Time    `json:"createdTime,omitempty" bson:"created_time,omitempty"`

	ConnectionID bson.ObjectId             `json:"connectionId,omitempty" bson:"connection_id,omitempty"`
	MessageID    *bson.ObjectId            `json:"messageId,omitempty" bson:"message_id,omitempty"`
	UserID       string                    `json:"userId,omitempty" bson:"user_id,omitempty"`
	Message      string                    `json:"message,omitempty" bson:"message,omitempty"`
	Type         enum.MessageTypeEnumValue `json:"type,omitempty" bson:"type,omitempty"`
	IsSuccess    *bool                     `json:"isSuccess,omitempty" bson:"is_success,omitempty"`
}

var MessageLogDB = &sdk.DBModel2{
	ColName:        "message_log",
	TemplateObject: &MessageLog{},
}

func InitMessageLogDB(dbSession *sdk.DBSession, dbName string) {
	MessageLogDB.DBName = dbName
	MessageLogDB.Init(dbSession)

	MessageLogDB.CreateIndex(mgo.Index{
		Key:        []string{"message_id"},
		Background: true,
		Name:       "message_id",
	})

	MessageLogDB.CreateIndex(mgo.Index{
		Key:        []string{"user_id"},
		Background: true,
		Name:       "user_id",
	})

	MessageLogDB.CreateIndex(mgo.Index{
		Key:        []string{"connection_id"},
		Background: true,
		Name:       "connection_id",
	})
}
