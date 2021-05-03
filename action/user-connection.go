package action

import (
	"strings"
	"time"

	"github.com/globalsign/mgo/bson"
	"github.com/mssola/user_agent"
	"gitlab.ghn.vn/common-projects/go-sdk/sdk"
	"gitlab.ghn.vn/internal-tools/message/model"
	"gitlab.ghn.vn/internal-tools/message/model/enum"
)

var USER_CON_ID = "userConId"
var USER_ID = "userId"
var CONNECTED_TIME = "connectedTime"

func AuthorizeV2(connId bson.ObjectId, userId, userAgent, IP string) *sdk.APIResponse {
	if userId == "" || len(userId) == 0 {
		return &sdk.APIResponse{
			Status:  sdk.APIStatus.Invalid,
			Message: "Invalid [SSOToken] field",
		}
	}

	//ssoResp := client.GenTokenSSOv2(ssotoken, userAgent, IP)
	//if ssoResp.Status == sdk.APIStatus.Ok {
	//	genTokenSSOv2 := ssoResp.Data.([]client.GenTokenSSOv2Data)[0]
	//	accessToken := genTokenSSOv2.AccessToken
	//
	//	result := client.VerifySSOv2Token(accessToken, userAgent, IP)
	//	if result.Status != sdk.APIStatus.Ok {
	//		return result
	//	}
	//	userInfo := result.Data.([]client.VerifyTokenSSOv2Data)[0]

	respUser := model.UserDB.QueryOne(model.User{UserID: userId})
	if respUser.Status == sdk.APIStatus.Ok {
		now := time.Now()
		updater := model.UserConnection{
			UserID:         userId,
			AuthorizedTime: &now,
		}

		if userAgent != "" {
			ua := user_agent.New(userAgent)
			updater.UserAgent = userAgent
			updater.DeviceOs = ua.OS()

			uaParts := strings.Split(updater.UserAgent, " ")
			lenP := len(uaParts)
			if lenP > 1 {
				deviceId := uaParts[lenP-2]
				if strings.HasPrefix(deviceId, "ID(") {
					updater.DeviceID = deviceId[3 : len(deviceId)-2]
				}
				appVersion := uaParts[lenP-1]
				if strings.HasPrefix(appVersion, "Version(") {
					updater.AppVersion = appVersion[8 : len(appVersion)-2]
				}
			}
		}

		model.UserConnectionDB.UpdateOne(model.UserConnection{ID: connId}, updater)

		return respUser
	}

	// send old messages that missed from last connected
	// if lastConnected != nil && updater.DeviceID != "" {
	//	go sendOldMessageFromLastConnected(
	//		connId,
	//		string(userInfo.UserID),
	//		updater.DeviceID,
	//		lastConnected)
	// }

	return &sdk.APIResponse{
		Status:  sdk.APIStatus.Forbidden,
		Message: "You don't have permission to do this action.",
	}

}

func GetActiveConnectionOfUser(userId string) *sdk.APIResponse {
	return model.UserConnectionDB.Query(model.UserConnection{
		UserID: userId,
		Status: enum.ConStatus.ACTIVE,
	}, 0, 10, false)
}

type ConQuery struct {
	Status          enum.ConStatusEnumValue `bson:"status"`
	LastMessageTime map[string]*time.Time   `bson:"last_message_time"`
}

func TryToCleanConnection() {

	connMap := model.GetConnMap()
	outputMsg := &model.WSOutputMessage{
		Topic: enum.Topic.CONNECTION,
		Content: &model.Content{
			Action: "PING",
		},
	}

	now := time.Now()
	for _, conn := range connMap {

		if conn.Attached[USER_ID] == nil && conn.Attached[CONNECTED_TIME] != nil {
			connectedTime := conn.Attached[CONNECTED_TIME].(time.Time)
			if connectedTime.Add(180 * time.Second).Before(now) {
				outputMsg := &model.WSOutputMessage{
					Topic: enum.Topic.CONNECTION,
					Content: &model.Content{
						Status:  "CLOSED",
						ErrCode: "Too long not authorized.",
					},
				}
				PushMessageToDevice(conn, outputMsg.String(), nil)
				conn.Close()
				continue
			}
		}

		err := PushMessageToDevice(conn, outputMsg.String(), nil)
		if err != nil {
			conn.Close()
		}
	}

	// remove old connection that's not ping
	limit := now.Add(-10 * time.Minute)
	conResp := model.UserConnectionDB.Query(&ConQuery{
		Status: enum.ConStatus.ACTIVE,
		LastMessageTime: map[string]*time.Time{
			"$lt": &limit,
		},
	}, 0, 100, false)
	if conResp.Status == sdk.APIStatus.Ok {
		conList := conResp.Data.([]*model.UserConnection)
		for _, oldcon := range conList {
			model.UserConnectionDB.UpdateOne(&model.UserConnection{
				ID: oldcon.ID,
			}, &model.UserConnection{
				Status: enum.ConStatus.CLOSED,
			})
		}
	}
}
