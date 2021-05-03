package model

import (
	"time"

	"github.com/globalsign/mgo/bson"
	"gitlab.ghn.vn/common-projects/go-sdk/sdk"
)

type MessageRoom struct {
	ID              bson.ObjectId `json:"id,omitempty" bson:"_id,omitempty"`
	LastUpdatedTime *time.Time    `json:"lastUpdatedTime,omitempty" bson:"last_updated_time,omitempty"`
	CreatedTime     *time.Time    `json:"createdTime,omitempty" bson:"created_time,omitempty"`

	UserId     string   `json:"userId,omitempty" bson:"user_id,omitempty"`
	ChatRoomId string   `json:"chatRoomId,omitempty" bson:"chat_room_id,omitempty"`
	Content    *Content `json:"content,omitempty" bson:"content,omitempty"`
}

var MessageRoomDB = &sdk.DBModel2{
	ColName:        "message_room",
	TemplateObject: &MessageRoom{},
}

func InitMessageRoomDB(dbSession *sdk.DBSession, dbName string) {
	MessageRoomDB.DBName = dbName
	MessageRoomDB.Init(dbSession)
}
