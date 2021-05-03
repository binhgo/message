package api

import (
	"bytes"
	"crypto/rsa"
	"encoding/json"
	"os"
	"strings"
	"time"

	"github.com/globalsign/mgo/bson"
	"gitlab.ghn.vn/common-projects/go-sdk/sdk"
	"gitlab.ghn.vn/common-projects/go-sdk/sdk/websocket"
	"gitlab.ghn.vn/internal-tools/message/action"
	"gitlab.ghn.vn/internal-tools/message/cip"
	"gitlab.ghn.vn/internal-tools/message/model"
	"gitlab.ghn.vn/internal-tools/message/model/enum"
)

func OnWSConnected(conn *websocket.Connection) {
	now := time.Now()
	hostname, _ := os.Hostname()
	version := os.Getenv("version")
	userCon := &model.UserConnection{
		Status:        enum.ConStatus.ACTIVE,
		UserID:        "UNDEFINED",
		ConnectedTime: &now,
		WSHost:        hostname,
		WSHostVersion: version,
		IP:            conn.GetIP(),
		UserAgent:     conn.GetUserAgent(),
		InServiceID:   conn.Id,
	}

	result := model.UserConnectionDB.UpsertOne(&model.UserConnection{
		UserID: "NONE",
	}, userCon)

	if result.Status == sdk.APIStatus.Ok {

		conn.Attached[action.USER_CON_ID] = result.Data.([]*model.UserConnection)[0].ID
		conn.Attached[action.CONNECTED_TIME] = now

		outputMsg := &model.WSOutputMessage{
			Topic: enum.Topic.CONNECTION,
			Content: &model.Content{
				Status:     result.Status,
				Data:       result.Data,
				ApiMessage: "Connected. Please authorize in 180 seconds.",
			},
		}
		action.PushMessageToDevice(conn, outputMsg.String(), nil)
	}
}

func OnWSMessage(conn *websocket.Connection, message string) {
	var msg model.WSInputMessage

	err := json.Unmarshal([]byte(message), &msg)
	if err != nil {
		outputMsg := &model.WSOutputMessage{
			Content: &model.Content{
				Status:  sdk.APIStatus.Error,
				ErrCode: "Parse JSON failed: " + err.Error(),
			},
		}
		action.PushMessageToDevice(conn, outputMsg.String(), nil)
		return
	}

	if len(msg.Topic) == 0 {
		return
	}

	switch msg.Topic {

	// case "PING_PONG":
	// 	casePingPong(conn, msg)

	case "AUTH_TEST":
		caseFakeAuth(conn, msg)
		break

	case enum.Topic.AUTHORIZATION:
		caseAuth(conn, msg)
		break

	// If topic file/image -> content message is url
	case enum.Topic.MESSAGE, enum.Topic.IMAGE, enum.Topic.FILE:
		caseMsg(conn, msg)
		break

	case enum.Topic.VIDEO_OFFER, enum.Topic.VIDEO_ANSWER, enum.Topic.NEW_ICE_CANDIDATE:
		caseSignaling(conn, msg)
		break
	}

	// Log
	log := model.MessageLog{
		Message:      message,
		ConnectionID: conn.Attached[action.USER_CON_ID].(bson.ObjectId),
		Type:         enum.MessageType.RECEIVE,
	}

	if str, ok := conn.Attached[action.USER_ID].(string); ok {
		log.UserID = str
	}

	go model.MessageLogDB.Create(log)
	// end: Log

}

func OnWSClose(conn *websocket.Connection, err error) {
	now := time.Now()

	model.UserConnectionDB.UpdateOne(&model.UserConnection{
		ID: conn.Attached[action.USER_CON_ID].(bson.ObjectId),
	}, &model.UserConnection{
		Status:           enum.ConStatus.CLOSED,
		DisconnectedTime: &now,
		ClosedMessage:    err.Error(),
	})
}

func caseAuth(conn *websocket.Connection, msg model.WSInputMessage) {
	var data model.Content
	byteArr, _ := json.Marshal(msg.Content)
	err := json.Unmarshal(byteArr, &data)

	if err != nil {
		outputMsg := &model.WSOutputMessage{
			Topic: enum.Topic.AUTHORIZATION,
			Content: &model.Content{
				Status:  sdk.APIStatus.Invalid,
				ErrCode: "Parse data error: " + err.Error(),
			},
		}

		action.PushMessageToDevice(conn, outputMsg.String(), nil)
	}

	if len(data.UserId) == 0 {
		outputMsg := &model.WSOutputMessage{
			Topic: enum.Topic.AUTHORIZATION,
			Content: &model.Content{
				Status:  sdk.APIStatus.Invalid,
				ErrCode: "Require [SSOToken] field.",
			},
		}

		action.PushMessageToDevice(conn, outputMsg.String(), nil)
	}

	authResult := action.AuthorizeV2(
		conn.Attached[action.USER_CON_ID].(bson.ObjectId),
		data.UserId,
		msg.UserAgent,
		// conn.GetUserAgent(),
		conn.GetIP())

	if authResult.Status == sdk.APIStatus.Ok {
		userInfo := authResult.Data.([]*model.User)[0]
		userId := userInfo.UserID
		conn.Attached[action.USER_ID] = userId

		// add user
		filter := &model.User{UserID: userId}
		queryRs := model.UserDB.QueryOne(filter)

		var user *model.User

		switch queryRs.Status {

		case sdk.APIStatus.Ok:
			user = queryRs.Data.([]*model.User)[0]

			if len(user.Key) == 0 && len(user.PubKey) == 0 {
				priKey, pubKey, err := cip.GenerateRsaKeyPairPem()
				if err != nil {
					outputMsg := &model.WSOutputMessage{
						Topic: enum.Topic.AUTHORIZATION,
						Content: &model.Content{
							Status:  sdk.APIStatus.Invalid,
							ErrCode: err.Error(),
						},
					}

					action.PushMessageToDevice(conn, outputMsg.String(), nil)
					break
				}

				updateRs := model.UserDB.UpsertOne(filter, model.User{
					Key:    priKey,
					PubKey: pubKey,
				})

				if updateRs.Status == sdk.APIStatus.Ok {
					user = updateRs.Data.([]*model.User)[0]
				}
			}

			break

		//case sdk.APIStatus.NotFound:
		//
		//	priKey, pubKey, err := cip.GenerateRsaKeyPairPem()
		//	if err != nil {
		//		outputMsg := &model.WSOutputMessage{
		//			Topic: enum.Topic.AUTHORIZATION,
		//			Content: &model.Content{
		//				Status:  sdk.APIStatus.Invalid,
		//				ErrCode: err.Error(),
		//			},
		//		}
		//
		//		action.PushMessageToDevice(conn, outputMsg.String(), nil)
		//		break
		//	}
		//
		//	createRs := model.UserDB.Create(model.User{
		//		UserID: userId,
		//		Name:   userInfo.Name,
		//		Phone:  userInfo.Phone,
		//		Key:    priKey,
		//		PubKey: pubKey,
		//	})
		//
		//	if createRs.Status == sdk.APIStatus.Ok {
		//		user = createRs.Data.([]*model.User)[0]
		//	}
		//
		//	break
		}

		outputMsg := &model.WSOutputMessage{
			Topic: enum.Topic.AUTHORIZATION,
			Content: &model.Content{
				Status:     sdk.APIStatus.Ok,
				ApiMessage: authResult.Message,
				User:       user,
			},
		}

		action.PushMessageToDevice(conn, outputMsg.String(), nil)

	} else {

		outputMsg := &model.WSOutputMessage{
			Topic: enum.Topic.AUTHORIZATION,
			Content: &model.Content{
				Status:     authResult.Status,
				ApiMessage: authResult.Message,
				ErrCode:    authResult.ErrorCode,
			},
		}

		action.PushMessageToDevice(conn, outputMsg.String(), nil)
	}
}

func EncryptMsg(msg string, key *rsa.PublicKey) []byte {
	if strings.Contains(msg, "http") {
		return []byte(msg)
	}

	bb := bytes.Buffer{}
	bb.WriteString(msg)

	encryptedByte := cip.EncryptData(key, bb.Bytes())

	return encryptedByte
}

func caseMsg(conn *websocket.Connection, msg model.WSInputMessage) {

	if msg.Content == nil {
		outputMsg := &model.WSOutputMessage{
			Content: &model.Content{
				Status:  sdk.APIStatus.Invalid,
				ErrCode: "Content is empty",
			},
		}
		action.PushMessageToDevice(conn, outputMsg.String(), nil)
		return
	}

	// Get room
	roomResp := model.RoomDB.QueryOne(model.Room{
		ID: bson.ObjectIdHex(msg.Content.ChatRoomID),
	})

	if roomResp.Status != sdk.APIStatus.Ok {
		outputMsg := &model.WSOutputMessage{
			Content: &model.Content{
				Status:  sdk.APIStatus.NotFound,
				ErrCode: "NotFound any match room.",
			},
		}
		action.PushMessageToDevice(conn, outputMsg.String(), nil)
		return
	}

	roomData := roomResp.Data.([]*model.Room)[0]

	// todo
	// content for sender, will be encrypted later by sender Public Key
	msg.Content.SenderChatMessage = msg.Content.ChatMessage

	// assign default value
	encryptedByte := []byte(msg.Content.ChatMessage)

	if roomData.Type == enum.RoomType.ONE {

		receiverId := action.GetReceiverId(roomData.UserIds, msg.Content.FromUserId)

		// encrypt
		userFilter := model.User{UserID: receiverId}
		queryRs := model.UserDB.QueryOne(userFilter)

		if queryRs.Status == sdk.APIStatus.Ok {
			receiver := queryRs.Data.([]*model.User)[0]
			pubKey, err := cip.PemToPublicKey(receiver.PubKey)
			if err == nil && pubKey != nil {
				encryptedByte = EncryptMsg(msg.Content.ChatMessage, pubKey)
			}
		}

	} else if roomData.Type == enum.RoomType.MANY {

		pubKey, err := cip.PemToPublicKey(roomData.PubKey)
		if err == nil && pubKey != nil {
			encryptedByte = EncryptMsg(msg.Content.ChatMessage, pubKey)
		}
	}

	msg.Content.ChatMessageByte = encryptedByte
	msg.Content.ChatMessage = string(encryptedByte)
	// end

	// create message in room
	msg.Content.Topic = string(msg.Topic)
	model.MessageRoomDB.Create(model.MessageRoom{
		UserId:     conn.Attached[action.USER_ID].(string),
		ChatRoomId: roomData.ID.Hex(),
		Content:    msg.Content,
	})

	// send message to client
	for _, userId := range roomData.UserIds {

		conQResult := model.UserConnectionDB.QueryOne(model.UserConnection{
			UserID: userId,
			Status: enum.ConStatus.ACTIVE,
		})

		if conQResult.Status != sdk.APIStatus.Ok {
			action.FirebaseMessageQueueDB.PushWithKeys(msg.Content, &[]string{msg.Content.FromUserId})
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
}

func caseFakeAuth(conn *websocket.Connection, msg model.WSInputMessage) {
	var data model.Content
	byteArr, _ := json.Marshal(msg.Content)
	_ = json.Unmarshal(byteArr, &data)

	updater := model.UserConnection{
		UserID: data.UserId,
	}
	conn.Attached[action.USER_ID] = data.UserId
	resp := model.UserConnectionDB.UpdateOne(model.UserConnection{ID: conn.Attached[action.USER_CON_ID].(bson.ObjectId)}, updater)

	outputMsg := &model.WSOutputMessage{
		Topic: "AUTH_TEST",
		Content: &model.Content{
			Status:     resp.Status,
			ApiMessage: resp.Message,
		},
	}

	action.PushMessageToDevice(conn, outputMsg.String(), nil)
}

// func casePingPong(conn *websocket.Connection, msg model.WSInputMessage) {
// 	action.CreateNewMessage(msg.Topic, msg.Content.FromUserId, msg.Content)
// }

func caseSignaling(conn *websocket.Connection, msg model.WSInputMessage) {
	if msg.Content == nil {
		outputMsg := &model.WSOutputMessage{
			Content: &model.Content{
				Status:  sdk.APIStatus.Invalid,
				ErrCode: "Content is empty",
			},
		}
		action.PushMessageToDevice(conn, outputMsg.String(), nil)
		return
	}

	// Get room
	roomResp := model.RoomDB.QueryOne(model.Room{
		ID: bson.ObjectIdHex(msg.Content.ChatRoomID),
	})

	if roomResp.Status != sdk.APIStatus.Ok {
		outputMsg := &model.WSOutputMessage{
			Content: &model.Content{
				Status:  sdk.APIStatus.NotFound,
				ErrCode: "NotFound any match room.",
			},
		}
		action.PushMessageToDevice(conn, outputMsg.String(), nil)
		return
	}

	roomData := roomResp.Data.([]*model.Room)[0]

	// send message to client
	for _, userId := range roomData.UserIds {

		conQResult := model.UserConnectionDB.QueryOne(model.UserConnection{
			UserID: userId,
			Status: enum.ConStatus.ACTIVE,
		})

		// if offer call to offline user -> noti for user
		if conQResult.Status != sdk.APIStatus.Ok && msg.Topic == enum.Topic.VIDEO_OFFER {
			action.FirebaseMessageQueueDB.PushWithKeys(msg.Content, &[]string{msg.Content.FromUserId})
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
}
