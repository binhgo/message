package action

import (
	"github.com/binhgo/go-sdk/sdk"
	"github.com/globalsign/mgo"
)

type DeviceRegistration struct {
	UserId         string `bson:"user_id,omitempty" json:"userId,omitempty"`
	RegistrationId string `bson:"registration_id,omitempty" json:"registrationId,omitempty"`
	DeviceName     string `bson:"device_name,omitempty" json:"deviceName,omitempty"`
}

var DeviceRegistrationDB = &sdk.DBModel2{
	ColName:        "firebase_registration",
	TemplateObject: &DeviceRegistration{},
}

func InitDeviceRegistrationDB(dbSession *sdk.DBSession, dbName string) {
	DeviceRegistrationDB.DBName = dbName
	DeviceRegistrationDB.Init(dbSession)

	DeviceRegistrationDB.CreateIndex(mgo.Index{
		Key:        []string{"user_id", "registration_id"},
		Background: true,
		Unique:     true,
	})

	DeviceRegistrationDB.CreateIndex(mgo.Index{
		Key:        []string{"user_id"},
		Background: true,
	})
}
