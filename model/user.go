package model

import (
	"github.com/binhgo/go-sdk/sdk"
	"github.com/globalsign/mgo/bson"
)

type User struct {
	ID bson.ObjectId `json:"id,omitempty" bson:"_id,omitempty"`

	UserID   string `json:"userId,omitempty" bson:"user_id,omitempty"`
	Password string `json:"password,omitempty" bson:"password,omitempty"`
	Name     string `json:"name,omitempty" bson:"name,omitempty"`
	Email    string `json:"email,omitempty" bson:"email,omitempty"`
	Phone    string `json:"phone,omitempty" bson:"phone,omitempty"`

	Avatar string `json:"avatar,omitempty" bson:"avatar,omitempty"`
	Key    string `json:"-" bson:"key,omitempty"`
	PubKey string `json:"pub_key,omitempty" bson:"pub_key,omitempty"`
}

var UserDB = &sdk.DBModel2{
	ColName:        "user",
	TemplateObject: &User{},
}

func InitUserDB(dbSession *sdk.DBSession, dbName string) {
	UserDB.DBName = dbName
	UserDB.Init(dbSession)
}
