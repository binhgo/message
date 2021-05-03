package model

import (
	"time"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"gitlab.ghn.vn/common-projects/go-sdk/sdk"
	"gitlab.ghn.vn/internal-tools/message/model/enum"
)

type Message struct {
	ID              bson.ObjectId `json:"id,omitempty" bson:"_id,omitempty"`
	LastUpdatedTime *time.Time    `json:"lastUpdatedTime,omitempty" bson:"last_updated_time,omitempty"`

	Topic          enum.TopicEnumValue         `json:"topic,omitempty" bson:"topic,omitempty"`
	ToUserId       string                      `json:"toUserId,omitempty" bson:"to_user_id,omitempty"`
	DeliveryStatus enum.MessageStatusEnumValue `json:"deliveryStatus,omitempty" bson:"delivery_status,omitempty"`
	Content        *Content                    `json:"content,omitempty" bson:"content,omitempty"`

	Targets    *[]*UserConnection `json:"targets,omitempty" bson:"targets,omitempty"`
	TargetSent *int               `json:"targetSent,omitempty" bson:"target_sent,omitempty"`
}

var MessageDB = &sdk.DBModel2{
	ColName:        "message",
	TemplateObject: &Message{},
}

func InitMessageDB(dbSession *sdk.DBSession, dbName string) {
	MessageDB.DBName = dbName
	MessageDB.Init(dbSession)

	MessageDB.CreateIndex(mgo.Index{
		Key:        []string{"to_user_id", "last_updated_time", "_id"},
		Background: true,
		Name:       "to_user_id_target_sent_last_updated_time",
		Sparse:     true,
	})

	MessageDB.CreateIndex(mgo.Index{
		Key:        []string{"to_user_id", "_id"},
		Background: true,
		Name:       "to_user_id",
		Sparse:     true,
	})

	MessageDB.CreateIndex(mgo.Index{
		Key:        []string{"to_user_id", "status", "_id"},
		Background: true,
		Sparse:     true,
		Name:       "to_user_id_status",
	})

	MessageDB.CreateIndex(mgo.Index{
		Key:        []string{"status", "_id"},
		Sparse:     true,
		Background: true,
		Name:       "status",
	})

	MessageDB.CreateIndex(mgo.Index{
		Key:        []string{"topic", "_id"},
		Sparse:     true,
		Background: true,
		Name:       "topic",
	})
}

type Content struct {
	Topic      string `json:"topic,omitempty" bson:"topic,omitempty"`
	Status     string `json:"status,omitempty" bson:"status,omitempty"`
	ApiMessage string `json:"apiMessage,omitempty" bson:"apiMessage,omitempty"`
	ErrCode    string `json:"errorCode,omitempty" bson:"errorCode,omitempty"`

	// encoded chat msg of receiver
	ChatMessage     string `json:"chatMessage,omitempty" bson:"chatMessage,omitempty"`
	ChatMessageByte []byte `json:"chatMessageByte,omitempty" bson:"chatMessageByte,omitempty"`

	// encoded chat msg of sender
	SenderChatMessage     string `json:"senderChatMessage,omitempty" bson:"senderChatMessage,omitempty"`
	SenderChatMessageByte []byte `json:"senderChatMessageByte,omitempty" bson:"senderChatMessageByte,omitempty"`

	// for signaling call
	SDP       interface{} `json:"sdp,omitempty" bson:"sdp,omitempty"`
	Candidate interface{} `json:"candidate,omitempty" bson:"candidate,omitempty"`

	Type        string      `json:"type,omitempty" bson:"type,omitempty"`
	ChatRoomID  string      `json:"chatRoomId,omitempty" bson:"chatRoomId,omitempty"`
	FromUserId  string      `json:"fromUserId,omitempty" bson:"fromUserId,omitempty"`
	//SSOToken    string      `json:"ssoToken,omitempty" bson:"ssoToken,omitempty"`
	UserId 		string 		`json:"userId,omitempty" bson:"user_id,omitempty"`
	Action      string      `json:"action,omitempty" bson:"action,omitempty"`
	Data        interface{} `json:"data,omitempty" bson:"data,omitempty"`
	User        *User       `json:"user,omitempty" bson:"user,omitempty"`
	IsVideoCall bool        `json:"isVideoCall,omitempty" bson:"is_video_call,omitempty"`

	LastUpdatedTime *time.Time `json:"lastUpdatedTime,omitempty" bson:"last_updated_time,omitempty"`
}
