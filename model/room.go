package model

import (
	"github.com/globalsign/mgo/bson"
	"gitlab.ghn.vn/common-projects/go-sdk/sdk"
	"gitlab.ghn.vn/internal-tools/message/model/enum"
)

type Room struct {
	ID bson.ObjectId `json:"id,omitempty" bson:"_id,omitempty"`

	UserIds    []string `json:"userIds,omitempty" bson:"user_ids,omitempty"`
	PinMessage *Content `json:"pinMessage,omitempty" bson:"pin_message,omitempty"`
	OffsetPin  int64    `json:"offsetPin,omitempty" bson:"offset_pin,omitempty"` // ->use for scroll more message near pin
	Name       string   `json:"name,omitempty" bson:"name,omitempty"`
	Avatar     string   `json:"avatar,omitempty" bson:"avatar,omitempty"`

	Type    enum.RoomTypeEnumValue `json:"roomType,omitempty" bson:"room_type,omitempty"`
	RoomKey string                 `json:"roomKey,omitempty" bson:"room_key,omitempty"`

	PriKey string `json:"-" bson:"pri_key,omitempty"`
	PubKey string `json:"pubKey,omitempty" bson:"pub_key,omitempty"`
}

var RoomDB = &sdk.DBModel2{
	ColName:        "room",
	TemplateObject: &Room{},
}

func InitRoomDB(dbSession *sdk.DBSession, dbName string) {
	RoomDB.DBName = dbName
	RoomDB.Init(dbSession)
}
