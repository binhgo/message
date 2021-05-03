package api

import (
	"github.com/binhgo/go-sdk/sdk"
	"github.com/globalsign/mgo/bson"

	"github.com/binhgo/message/action"
	"github.com/binhgo/message/model"
	"github.com/binhgo/message/model/enum"
)

// create message
// if have fromUserId -> don't send to this user
func CreateMessage(req sdk.APIRequest, res sdk.APIResponder) error {
	var msg model.WSInputMessage

	if err := req.GetContent(&msg); err != nil {
		return res.Respond(&sdk.APIResponse{
			Status:    sdk.APIStatus.Invalid,
			ErrorCode: "FAIL_TO_PARSE_JSON",
			Message:   "Fail to parse JSON, please check the format. " + err.Error(),
		})
	}

	if msg.Content == nil || len(msg.Topic) == 0 {
		return res.Respond(&sdk.APIResponse{
			Status:    sdk.APIStatus.Invalid,
			ErrorCode: "INVALID_DATA",
			Message:   "INVALID data body.",
		})
	}

	// Get room
	roomResp := model.RoomDB.QueryOne(model.Room{
		ID: bson.ObjectIdHex(msg.Content.ChatRoomID),
	})

	if roomResp.Status != sdk.APIStatus.Ok {
		return res.Respond(roomResp)
	}

	roomData := roomResp.Data.([]*model.Room)[0]

	// send message to client
	for _, userId := range roomData.UserIds {

		if len(msg.Content.FromUserId) > 0 && userId == msg.Content.FromUserId {
			continue
		}
		conQResult := model.UserConnectionDB.QueryOne(model.UserConnection{
			UserID: userId,
			Status: enum.ConStatus.ACTIVE,
		})

		// if offer call to offline user -> noti for user
		if conQResult.Status != sdk.APIStatus.Ok {
			if msg.Topic == enum.Topic.VIDEO_OFFER {
				action.FirebaseMessageQueueDB.PushWithKeys(msg.Content, &[]string{msg.Content.FromUserId})
			}

			continue
		}

		uConn := conQResult.Data.([]*model.UserConnection)[0]
		go func(con *model.UserConnection, msg model.WSInputMessage) {

			// Push message to user
			if con.Status == enum.ConStatus.ACTIVE {
				action.CreateNewMessage(msg.Topic, con.UserID, msg.Content)
			}

		}(uConn, msg)

	}

	return res.Respond(&sdk.APIResponse{
		Status:  sdk.APIStatus.Ok,
		Message: "Create message success.",
	})
}
