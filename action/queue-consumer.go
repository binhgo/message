package action

import (
	"errors"

	"github.com/binhgo/go-sdk/sdk"
	"github.com/globalsign/mgo/bson"

	"github.com/binhgo/message/cip"
	"github.com/binhgo/message/model"
	"github.com/binhgo/message/model/enum"
)

func QueueConsume(item *sdk.SortedQueueItem) error {

	byteArr, err := bson.Marshal(item.Data)
	if err != nil {
		return err
	}

	var qItem model.MessageQueueItem
	err = bson.Unmarshal(byteArr, &qItem)
	if err != nil {
		return err
	}

	switch qItem.Category {

	case enum.PUSH_MESSAGE:

		queryRs := model.MessageDB.QueryOne(&model.Message{
			ID: bson.ObjectIdHex(qItem.Data),
		})

		if queryRs.Status == sdk.APIStatus.Ok {
			msg := queryRs.Data.([]*model.Message)[0]

			if msg.Topic == enum.Topic.VIDEO_OFFER || msg.Topic == enum.Topic.VIDEO_ANSWER ||
				msg.Topic == enum.Topic.NEW_ICE_CANDIDATE || msg.Topic == enum.Topic.HANG_UP {
				// do nothing
			} else {
				decrypt(msg)
			}

			conn := model.GetConnById(qItem.ConnectionLocalID)
			if conn == nil {
				FirebaseMessageQueueDB.PushWithKeys(msg.Content, &[]string{msg.Content.FromUserId})
				return nil
			} else {

				msg.Content.LastUpdatedTime = msg.LastUpdatedTime
				outputMsg := &model.WSOutputMessage{
					Topic:   msg.Topic,
					Content: msg.Content,
				}

				err := PushMessageToDevice(conn, outputMsg.String(), nil)
				if err == nil {
					go logMsg(msg)
				} else {
					return err
				}
			}

		} else if queryRs.Status == sdk.APIStatus.NotFound {
			return nil
		} else {
			return errors.New("ERROR QUERY MESSAGE: " + queryRs.Message)
		}

		break
	}

	return nil
}

func decrypt(msg *model.Message) {

	// Get room
	roomRs := model.RoomDB.QueryOne(model.Room{
		ID: bson.ObjectIdHex(msg.Content.ChatRoomID),
	})

	// assign default value
	decryptedMsg := "WTH: Error during decrypting message."

	if roomRs.Status == sdk.APIStatus.Ok {
		room := roomRs.Data.([]*model.Room)[0]
		if room.Type == enum.RoomType.MANY {

			priKey, err := cip.PemToPrivateKey(room.PriKey)
			if err == nil && priKey != nil {
				decryptedMsg = DecryptMsg(msg.Content.ChatMessageByte, priKey)
			}

		} else if room.Type == enum.RoomType.ONE {

			receiverId := GetReceiverId(room.UserIds, msg.Content.FromUserId)

			// de-cr
			userFilter := model.User{UserID: receiverId}
			queryRs := model.UserDB.QueryOne(userFilter)

			if queryRs.Status == sdk.APIStatus.Ok {
				receiver := queryRs.Data.([]*model.User)[0]
				priKey, err := cip.PemToPrivateKey(receiver.Key)
				if err == nil && priKey != nil {
					decryptedMsg = DecryptMsg(msg.Content.ChatMessageByte, priKey)
				}
			}
		}

	} else {
		decryptedMsg = "WTH: Room not found: " + roomRs.Message
	}

	msg.Content.ChatMessage = decryptedMsg
}

func logMsg(msg *model.Message) {

	filter := &model.Message{ID: msg.ID}
	qResult := model.MessageDB.IncreOne(filter, "target_sent", 1)

	if qResult.Status == sdk.APIStatus.Ok {

		msg = qResult.Data.([]*model.Message)[0]

		if msg.TargetSent != nil && *msg.TargetSent == len(*msg.Targets) {
			updater := &model.Message{DeliveryStatus: enum.MessageStatus.DELIVERED}
			model.MessageDB.UpdateOne(filter, updater)
		}
	}
}
