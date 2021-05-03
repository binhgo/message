package action

import (
	"time"

	"github.com/binhgo/go-sdk/sdk"
	"github.com/binhgo/go-sdk/sdk/websocket"
	"github.com/globalsign/mgo/bson"

	"github.com/binhgo/message/model"
	"github.com/binhgo/message/model/enum"
)

var supportedTopic = []enum.TopicEnumValue{
	enum.Topic.AUTHORIZATION,
	enum.Topic.ROOM_NOTIFY,
	enum.Topic.MESSAGE,
	enum.Topic.IMAGE,
	enum.Topic.FILE,
	enum.Topic.NEW_ICE_CANDIDATE,
	enum.Topic.VIDEO_ANSWER,
	enum.Topic.VIDEO_OFFER,
	enum.Topic.HANG_UP,
}

func IsValidTopic(input enum.TopicEnumValue) bool {
	for _, topic := range supportedTopic {
		if input == topic {
			return true
		}
	}
	return false
}

func CreateNewMessage(topic enum.TopicEnumValue, toUserId string, content *model.Content) *sdk.APIResponse {

	if content == nil {
		return &sdk.APIResponse{
			Status:    sdk.APIStatus.Invalid,
			Message:   "[content] must not empty.",
			ErrorCode: "EMPTY_CONTENT",
		}
	}

	if topic == "" || !IsValidTopic(topic) {
		return &sdk.APIResponse{
			Status:    sdk.APIStatus.Invalid,
			Message:   "[topic] is not accepted.",
			ErrorCode: "INVALID_TOPIC",
		}
	}

	if toUserId == "" {
		return &sdk.APIResponse{
			Status:    sdk.APIStatus.Invalid,
			Message:   "[toUserId] must not empty",
			ErrorCode: "INVALID_USER",
		}
	}

	zero := 0

	msg := &model.Message{
		DeliveryStatus: enum.MessageStatus.DELIVERING,
		Topic:          topic,
		Content:        content,
		ToUserId:       toUserId,
		Targets:        &[]*model.UserConnection{},
		TargetSent:     &zero,
	}

	// scan target devices of driver
	conQResult := model.UserConnectionDB.Query(model.UserConnection{
		UserID: toUserId,
		Status: enum.ConStatus.ACTIVE,
	}, 0, 20, false)

	if conQResult.Status == sdk.APIStatus.Ok {
		conList, ok := conQResult.Data.([]*model.UserConnection)
		if ok && len(conList) > 0 {
			for _, con := range conList {
				con.DeliverType = enum.DeliverType.DIRECT
			}
			msg.Targets = &conList
		}
	}

	if msg.Targets == nil || len(*msg.Targets) == 0 {
		msg.DeliveryStatus = enum.MessageStatus.DELIVERED
	}

	// use upsert to get created object _id value
	createResult := model.MessageDB.UpsertOne(&model.Message{
		Topic: enum.Topic.NONE,
	}, msg)

	// push to queue
	if createResult.Status == sdk.APIStatus.Ok {

		conList := *msg.Targets
		msg = createResult.Data.([]*model.Message)[0]

		for _, con := range conList {

			queueItem := model.MessageQueueItem{
				Category:          enum.PUSH_MESSAGE,
				ConnectionLocalID: con.InServiceID,
				Data:              msg.ID.Hex(),
			}

			logKey := &[]string{
				"TO_USER/" + toUserId,
				"MSG/" + msg.ID.Hex(),
			}

			model.MessageQueueDB.PushWithKeysAndTopic(queueItem, toUserId, logKey, con.WSHost)
		}
	}

	return createResult
}

func PushMessageToDevice(conn *websocket.Connection, message string, messageId *bson.ObjectId) error {
	if !conn.IsActive() {
		return &sdk.Error{Type: "DISCONNECTED", Message: "Connection closed."}
	}
	err := conn.Send(message)
	isSuccess := err == nil
	conId := conn.Attached["userConId"].(bson.ObjectId)
	log := model.MessageLog{
		Message:      message,
		ConnectionID: conId,
		IsSuccess:    &isSuccess,
		MessageID:    messageId,
		Type:         enum.MessageType.PUSH,
	}

	if str, ok := conn.Attached["userId"].(string); ok {
		log.UserID = str
	}
	go model.MessageLogDB.Create(log)

	if isSuccess {
		now := time.Now()
		model.UserConnectionDB.UpdateOne(model.UserConnection{
			ID: conId,
		}, model.UserConnection{
			LastMessageTime: &now,
		})
	}

	return err
}
