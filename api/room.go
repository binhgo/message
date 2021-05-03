package api

import (
	"sort"
	"strings"

	"github.com/globalsign/mgo/bson"
	"gitlab.ghn.vn/common-projects/go-sdk/sdk"
	"gitlab.ghn.vn/internal-tools/message/action"
	"gitlab.ghn.vn/internal-tools/message/cip"
	"gitlab.ghn.vn/internal-tools/message/model"
	"gitlab.ghn.vn/internal-tools/message/model/enum"
)

func CreateRoom(req sdk.APIRequest, res sdk.APIResponder) error {

	var input model.Room

	if err := req.GetContent(&input); err != nil {
		return res.Respond(&sdk.APIResponse{
			Status:    sdk.APIStatus.Invalid,
			ErrorCode: "FAIL_TO_PARSE_JSON",
			Message:   "Fail to parse JSON, please check the format. " + err.Error(),
		})
	}

	if len(input.UserIds) == 2 {
		userIds := input.UserIds
		sort.Strings(userIds)

		roomKey := strings.Join(userIds, ".")
		check := model.RoomDB.QueryOne(model.Room{
			RoomKey: roomKey,
			Type:    enum.RoomType.ONE,
		})
		if check.Status == sdk.APIStatus.Ok {
			return res.Respond(&sdk.APIResponse{
				Status:  sdk.APIStatus.Existed,
				Data:    check.Data,
				Message: "This room already created.",
			})
		}

		input.RoomKey = roomKey
		input.Type = enum.RoomType.ONE

	} else {

		input.Type = enum.RoomType.MANY

		// gen key
		priKey, pubKey, err := cip.GenerateRsaKeyPairPem()
		if err != nil {
			return res.Respond(&sdk.APIResponse{
				Status:    sdk.APIStatus.Invalid,
				ErrorCode: "ERR_GEN_KEY_PAIR",
				Message:   "Error gen key pair",
			})
		}

		input.PriKey = priKey
		input.PubKey = pubKey
	}

	resp := model.RoomDB.Create(input)
	if resp.Status == sdk.APIStatus.Ok {
		roomResp := model.RoomDB.QueryOne(input)
		if roomResp.Status == sdk.APIStatus.Ok {
			roomInfo := roomResp.Data.([]*model.Room)[0]

			// send message to user
			go func(room *model.Room, userIds []string) {

				for _, userId := range userIds {

					content := &model.Content{
						Type: "NEW_ROOM",
						Data: room,
					}

					action.CreateNewMessage(enum.Topic.ROOM_NOTIFY, userId, content)

				}

			}(roomInfo, input.UserIds)

			return res.Respond(roomResp)
		}
	}

	return res.Respond(&sdk.APIResponse{
		Status:  sdk.APIStatus.Invalid,
		Message: "ERROR: Something wrong when create room.",
	})
}

func GetRoom(req sdk.APIRequest, res sdk.APIResponder) error {
	var userId = req.GetParam("userId")
	if userId == "" || len(userId) == 0 {
		return res.Respond(&sdk.APIResponse{
			Status:    sdk.APIStatus.Invalid,
			ErrorCode: "INVALID_PARAM",
			Message:   "Invalid [userId] param.",
		})
	}

	var offset = sdk.ParseInt(req.GetParam("offset"), 0)
	var limit = sdk.ParseInt(req.GetParam("limit"), 10)
	var reverse = req.GetParam("reverse") == "true"
	// var getTotal = req.GetParam("getTotal") == "true"

	filter := bson.M{
		"user_ids": userId,
	}

	return res.Respond(model.RoomDB.Query(filter, offset, limit, reverse))
}

func GetRoomInfo(req sdk.APIRequest, res sdk.APIResponder) error {
	var chatRoomId = req.GetParam("chatRoomId")
	if chatRoomId == "" || len(chatRoomId) == 0 {
		return res.Respond(&sdk.APIResponse{
			Status:    sdk.APIStatus.Invalid,
			ErrorCode: "INVALID_PARAM",
			Message:   "Invalid [chatRoomId] param.",
		})
	}

	filter := model.Room{
		ID: bson.ObjectIdHex(chatRoomId),
	}

	return res.Respond(model.RoomDB.QueryOne(filter))
}

func UpdateRoomInfoPUT(req sdk.APIRequest, res sdk.APIResponder) error {

	var input model.Room
	if err := req.GetContent(&input); err != nil {
		return res.Respond(&sdk.APIResponse{
			Status:    sdk.APIStatus.Invalid,
			ErrorCode: "FAIL_TO_PARSE_JSON",
			Message:   "Fail to parse JSON, please check the format. " + err.Error(),
		})
	}

	if len(input.ID) == 0 {
		return res.Respond(&sdk.APIResponse{
			Status:    sdk.APIStatus.Invalid,
			ErrorCode: "INVALID_ID",
			Message:   "Invalid [ChatroomID] parameter. ",
		})
	}

	filter := model.Room{ID: input.ID}
	resp := model.RoomDB.QueryOne(filter)
	if resp.Status != sdk.APIStatus.Ok {
		return res.Respond(resp)
	}

	updater := model.Room{}
	// update pin
	if input.PinMessage != nil && input.OffsetPin >= 0 {
		updater.PinMessage = input.PinMessage
		updater.OffsetPin = input.OffsetPin
	}
	// update name
	if input.Name != "" || len(input.Name) > 0 {
		updater.Name = input.Name
	}
	// update avatar
	if input.Avatar != "" || len(input.Avatar) > 0 {
		updater.Avatar = input.Avatar
	}

	return res.Respond(model.RoomDB.UpdateOne(filter, updater))
}

// POST add pin message
func PinMessage(req sdk.APIRequest, res sdk.APIResponder) error {
	var input model.Room

	if err := req.GetContent(&input); err != nil {
		return res.Respond(&sdk.APIResponse{
			Status:    sdk.APIStatus.Invalid,
			ErrorCode: "FAIL_TO_PARSE_JSON",
			Message:   "Fail to parse JSON, please check the format. " + err.Error(),
		})
	}

	if input.PinMessage == nil || len(input.ID) == 0 || input.OffsetPin < 0 {
		return res.Respond(&sdk.APIResponse{
			Status:    sdk.APIStatus.Invalid,
			ErrorCode: "INVALID_PARAM",
			Message:   "Invalid [PinMessage]/[ChatRoomID]/[OffsetPin] parameter.",
		})
	}

	filter := model.Room{ID: input.ID}
	resp := model.RoomDB.QueryOne(filter)
	if resp.Status != sdk.APIStatus.Ok {
		return res.Respond(resp)
	}

	updater := model.Room{
		PinMessage: input.PinMessage,
		OffsetPin:  input.OffsetPin,
	}

	return res.Respond(model.RoomDB.UpdateOne(filter, updater))
}

// Get click pin message
func GetPinMessage(req sdk.APIRequest, res sdk.APIResponder) error {

	chatRoomId := req.GetParam("chatRoomId")
	if len(chatRoomId) == 0 {
		return res.Respond(&sdk.APIResponse{
			Status:    sdk.APIStatus.Invalid,
			ErrorCode: "INVALID_PARAM",
			Message:   "Invalid [ChatRoomID] parameter.",
		})
	}

	resp := model.RoomDB.QueryOne(model.Room{ID: bson.ObjectIdHex(chatRoomId)})
	if resp.Status != sdk.APIStatus.Ok {
		return res.Respond(resp)
	}

	roomInfo := resp.Data.([]*model.Room)[0]
	return res.Respond(model.MessageRoomDB.Query(model.MessageRoom{
		Content:    roomInfo.PinMessage,
		ChatRoomId: chatRoomId,
	}, int(roomInfo.OffsetPin), 50, false))

}

func AddUserToRoomPUT(req sdk.APIRequest, res sdk.APIResponder) error {
	var input model.Room

	if err := req.GetContent(&input); err != nil {
		return res.Respond(&sdk.APIResponse{
			Status:    sdk.APIStatus.Invalid,
			ErrorCode: "FAIL_TO_PARSE_JSON",
			Message:   "Fail to parse JSON, please check the format. " + err.Error(),
		})
	}

	if len(input.ID) == 0 || len(input.UserIds) == 0 {
		return res.Respond(&sdk.APIResponse{
			Status:    sdk.APIStatus.Invalid,
			ErrorCode: "INVALID_ID",
			Message:   "Invalid chat room ID.",
		})
	}

	filter := model.Room{ID: input.ID}
	resp := model.RoomDB.QueryOne(filter)
	if resp.Status != sdk.APIStatus.Ok {
		return res.Respond(resp)
	}

	roomInfo := resp.Data.([]*model.Room)[0]
	var updater model.Room
	updater.UserIds = roomInfo.UserIds
	for _, userId := range input.UserIds {
		if action.IsExist(roomInfo.UserIds, userId) {
			continue
		}

		updater.UserIds = append(updater.UserIds, userId)
	}

	roomKey := strings.Join(updater.UserIds, ".")
	// update if room 1:1 to room many
	if len(roomInfo.UserIds) == 2 {
		updater.Type = enum.RoomType.MANY
		updater.RoomKey = roomKey
	}

	resp = model.RoomDB.UpdateOne(filter, updater)
	if resp.Status == sdk.APIStatus.Ok {
		roomIf := resp.Data.([]*model.Room)[0]
		go func(room, existedRoom *model.Room, userIds []string) {

			for _, userId := range userIds {
				if action.IsExist(existedRoom.UserIds, userId) {
					continue
				}

				content := &model.Content{
					Type: "NEW_ROOM",
					Data: room,
				}

				action.CreateNewMessage(enum.Topic.ROOM_NOTIFY, userId, content)
			}

		}(roomIf, roomInfo, input.UserIds)

	}

	return res.Respond(resp)
}

// one user per remove
func RemoveUserToRoomPUT(req sdk.APIRequest, res sdk.APIResponder) error {
	var input model.Room

	if err := req.GetContent(&input); err != nil {
		return res.Respond(&sdk.APIResponse{
			Status:    sdk.APIStatus.Invalid,
			ErrorCode: "FAIL_TO_PARSE_JSON",
			Message:   "Fail to parse JSON, please check the format. " + err.Error(),
		})
	}

	if len(input.ID) == 0 || len(input.UserIds) == 0 {
		return res.Respond(&sdk.APIResponse{
			Status:    sdk.APIStatus.Invalid,
			ErrorCode: "INVALID_ID",
			Message:   "Invalid chat room ID.",
		})
	}

	filter := model.Room{ID: input.ID}
	resp := model.RoomDB.QueryOne(filter)
	if resp.Status != sdk.APIStatus.Ok {
		return res.Respond(resp)
	}

	roomInfo := resp.Data.([]*model.Room)[0]

	// remove room when only one user in room
	if len(roomInfo.UserIds) == 2 {
		return res.Respond(model.RoomDB.Delete(filter))
	}

	var updater model.Room
	newUserId := action.RemoveUser(roomInfo.UserIds, input.UserIds[0])
	if newUserId == nil {
		return res.Respond(&sdk.APIResponse{
			Status:    sdk.APIStatus.Error,
			ErrorCode: "ERROR_INPUT",
			Message:   "ERROR: User not in room.",
		})
	}
	updater.UserIds = action.RemoveUser(roomInfo.UserIds, input.UserIds[0])

	// update type and key when room 1:1
	if len(updater.UserIds) == 2 {
		updater.Type = enum.RoomType.ONE
		userIds := updater.UserIds
		sort.Strings(userIds)

		roomKey := strings.Join(userIds, ".")
		updater.RoomKey = roomKey
	}

	updateResp := model.RoomDB.UpdateOne(filter, updater)
	if updateResp.Status == sdk.APIStatus.Ok {

		// remove room
		go func(room *model.Room, userId string) {

			content := &model.Content{
				Type: "REMOVE_ROOM",
				Data: room,
			}

			action.CreateNewMessage(enum.Topic.ROOM_NOTIFY, userId, content)

		}(roomInfo, input.UserIds[0])

	}

	return res.Respond(updateResp)
}
