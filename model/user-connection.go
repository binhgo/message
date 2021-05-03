package model

import (
	"os"
	"time"

	"github.com/binhgo/go-sdk/sdk"
	"github.com/binhgo/go-sdk/sdk/websocket"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"

	"github.com/binhgo/message/model/enum"
)

type UserConnection struct {
	ID              bson.ObjectId `json:"id,omitempty" bson:"_id,omitempty"`
	LastUpdatedTime *time.Time    `json:"lastUpdatedTime,omitempty" bson:"last_updated_time,omitempty"`
	CreatedTime     *time.Time    `json:"createdTime,omitempty" bson:"created_time,omitempty"`

	Status        enum.ConStatusEnumValue `json:"status,omitempty" bson:"status,omitempty"`
	WSHost        string                  `json:"wsHost,omitempty" bson:"ws_host,omitempty"`
	WSHostVersion string                  `json:"wsHostVersion,omitempty" bson:"ws_host_version,omitempty"`
	DeviceOs      string                  `json:"deviceOs,omitempty" bson:"device_os,omitempty"`
	AppVersion    string                  `json:"appVersion,omitempty" bson:"app_version,omitempty"`
	DeviceID      string                  `json:"deviceId,omitempty" bson:"device_id,omitempty"`
	IP            string                  `json:"ip,omitempty" bson:"ip,omitempty"`
	UserID        string                  `json:"userId,omitempty" bson:"user_id,omitempty"`
	UserAgent     string                  `json:"userAgent,omitempty" bson:"user_agent,omitempty"`
	InServiceID   int                     `json:"inServiceId,omitempty" bson:"in_service_id,omitempty"`
	ClosedMessage string                  `json:"closedMessage,omitempty" bson:"closed_message,omitempty"`

	ConnectedTime    *time.Time `json:"connectedTime,omitempty" bson:"connected_time,omitempty"`
	AuthorizedTime   *time.Time `json:"authorizedTime,omitempty" bson:"authorized_time,omitempty"`
	DisconnectedTime *time.Time `json:"disconnected_time,omitempty" bson:"disconnected_time,omitempty"`
	LastMessageTime  *time.Time `json:"lastMessageTime,omitempty" bson:"last_message_time,omitempty"`

	// used in message.Targets
	DeliverType enum.DeliverTypeEnumValue `json:"deliverType,omitempty" bson:"deliver_type,omitempty"`
}

var chatWSS *websocket.WSServer
var UserConnectionDB = &sdk.DBModel2{
	ColName:        "user_connection",
	TemplateObject: &UserConnection{},
}

type NotEqualStr struct {
	NeStr string `bson:"$ne,omitempty"`
}

type OldConnectionQuery struct {
	WSHostVersion NotEqualStr             `bson:"ws_host_version,omitempty"`
	Status        enum.ConStatusEnumValue `bson:"status,omitempty"`
}

func SetWebSocket(wss *websocket.WSServer) {
	chatWSS = wss
}

func GetConnById(id int) *websocket.Connection {
	return chatWSS.GetRoute("/").GetConnection(id)
}

func GetConnMap() map[int]*websocket.Connection {
	return chatWSS.GetRoute("/").GetConnectionMap()
}

func InitUserConDB(dbSession *sdk.DBSession, dbName string) {
	UserConnectionDB.DBName = dbName
	UserConnectionDB.Init(dbSession)

	UserConnectionDB.CreateIndex(mgo.Index{
		Key:        []string{"status", "last_message_time", "_id"},
		Background: true,
		Name:       "status_last_message_time",
	})

	UserConnectionDB.CreateIndex(mgo.Index{
		Key:        []string{"user_id", "status", "_id"},
		Background: true,
		Name:       "user_id_status",
	})

	UserConnectionDB.CreateIndex(mgo.Index{
		Key:        []string{"user_id", "status", "device_id"},
		Background: true,
		Name:       "user_id_device_id",
	})

	UserConnectionDB.CreateIndex(mgo.Index{
		Key:        []string{"ws_host", "status", "_id"},
		Background: true,
		Name:       "ws_host_status",
	})

	UserConnectionDB.CreateIndex(mgo.Index{
		Key:        []string{"ws_host_version", "status"},
		Background: true,
		Name:       "ws_host_version_status",
	})

	now := time.Now()
	UserConnectionDB.Update(&OldConnectionQuery{
		WSHostVersion: NotEqualStr{
			NeStr: os.Getenv("version"),
		},
		Status: enum.ConStatus.ACTIVE,
	}, &UserConnection{
		Status:           enum.ConStatus.CLOSED,
		DisconnectedTime: &now,
	})
}
