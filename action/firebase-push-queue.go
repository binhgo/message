package action

import (
	"github.com/globalsign/mgo/bson"
	"gitlab.ghn.vn/common-projects/go-sdk/sdk"
	"gitlab.ghn.vn/internal-tools/message/model"
)

//
// type FirebaseMessageQueueItem struct {
// 	Data       string             `bson:"data"`
// 	DeviceInfo DeviceRegistration `bson:"device_info"`
// }

var FirebaseMessageQueueDB = sdk.DBQueue2{
	ColName: "firebase_message_queue",
}

func InitFirebaseMessageQueueDB(dbSession *sdk.DBSession, dbName string) {
	FirebaseMessageQueueDB.Init(dbSession, dbName)
	FirebaseMessageQueueDB.StartConsumer(FirebaseMessageQueueDBConsumer, 5)
}

func FirebaseMessageQueueDBConsumer(item *sdk.QueueItem) error {

	// parse
	byteArr, err := bson.Marshal(item.Data)
	if err != nil {
		return err
	}

	var content model.Content
	err = bson.Unmarshal(byteArr, &content)
	if err != nil {
		return err
	}

	fbClient, err := GetDefaultFirebaseClient()
	if err != nil {
		return err
	}

	// query room
	roomRs := model.RoomDB.QueryOne(model.Room{
		ID: bson.ObjectIdHex(content.ChatRoomID),
	})

	if roomRs.Status == sdk.APIStatus.Ok {
		room := roomRs.Data.([]*model.Room)[0]
		for _, v := range room.UserIds {
			filter := DeviceRegistration{
				UserId: v,
			}

			queryRs := DeviceRegistrationDB.Query(filter, 0, 100, true)
			if queryRs.Status == sdk.APIStatus.Ok {
				devices := queryRs.Data.([]*DeviceRegistration)

				for _, d := range devices {
					msg := []FirebaseMessages{
						fbClient.WithToken(d.RegistrationId),
						fbClient.WithData(map[string]interface{}{
							"message": content.ChatMessage,
						}),
					}

					fbClient.Send(msg...)
				}
			}
		}
	}

	return nil
}
