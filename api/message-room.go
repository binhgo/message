package api

import (
	"github.com/globalsign/mgo/bson"
	"gitlab.ghn.vn/common-projects/go-sdk/sdk"
	"gitlab.ghn.vn/internal-tools/message/action"
	"gitlab.ghn.vn/internal-tools/message/cip"
	"gitlab.ghn.vn/internal-tools/message/model"
	"gitlab.ghn.vn/internal-tools/message/model/enum"
)

func GetMsgRoom(req sdk.APIRequest, res sdk.APIResponder) error {
	var chatRoomId = req.GetParam("chatRoomId")
	if chatRoomId == "" || len(chatRoomId) == 0 {
		return res.Respond(&sdk.APIResponse{
			Status:    sdk.APIStatus.Invalid,
			ErrorCode: "INVALID_PARAM",
			Message:   "Invalid [chatRoomId] param.",
		})
	}

	var offset = sdk.ParseInt(req.GetParam("offset"), 0)
	var limit = sdk.ParseInt(req.GetParam("limit"), 50)
	var reverse = req.GetParam("reverse") == "true"
	// var getTotal = req.GetParam("getTotal") == "true"

	filter := model.MessageRoom{
		ChatRoomId: chatRoomId,
	}

	queryRs := model.MessageRoomDB.Query(filter, offset, limit, reverse)
	if queryRs.Status == sdk.APIStatus.Ok {

		var results []*model.MessageRoom

		msgs := queryRs.Data.([]*model.MessageRoom)
		for _, v := range msgs {
			decrypt(v)
			results = append(results, v)
		}

		return res.Respond(&sdk.APIResponse{
			Status:  sdk.APIStatus.Ok,
			Data:    results,
			Message: "ok",
		})
	}

	return res.Respond(queryRs)
}

func decrypt(msg *model.MessageRoom) {

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
				decryptedMsg = action.DecryptMsg(msg.Content.ChatMessageByte, priKey)
			}

		} else if room.Type == enum.RoomType.ONE {

			receiverId := action.GetReceiverId(room.UserIds, msg.Content.FromUserId)

			// de-cr
			userFilter := model.User{UserID: receiverId}
			queryRs := model.UserDB.QueryOne(userFilter)

			if queryRs.Status == sdk.APIStatus.Ok {
				receiver := queryRs.Data.([]*model.User)[0]
				priKey, err := cip.PemToPrivateKey(receiver.Key)
				if err == nil && priKey != nil {
					decryptedMsg = action.DecryptMsg(msg.Content.ChatMessageByte, priKey)
				}
			}
		}

	} else {
		decryptedMsg = "WTH: Room not found: " + roomRs.Message
	}

	msg.Content.ChatMessage = decryptedMsg
}
