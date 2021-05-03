package api

import (
	"bytes"
	"errors"

	"github.com/binhgo/go-sdk/sdk"

	"github.com/binhgo/message/action"
)

func RegisterNewFirebaseDevice(req sdk.APIRequest, res sdk.APIResponder) error {

	var input action.DeviceRegistration
	if err := req.GetContent(&input); err != nil {
		return res.Respond(&sdk.APIResponse{
			Status:    sdk.APIStatus.Invalid,
			ErrorCode: "FAIL_TO_PARSE_JSON",
			Message:   "Fail to parse JSON, please check the format. " + err.Error(),
		})
	}

	if len(input.UserId) == 0 || len(input.RegistrationId) == 0 || len(input.DeviceName) == 0 {
		validateMsg := bytes.Buffer{}
		validateMsg.WriteString("require [UserId] [RegistrationId] [DeviceName]")

		return res.Respond(&sdk.APIResponse{
			Status:    sdk.APIStatus.Invalid,
			Message:   validateMsg.String(),
			ErrorCode: "INVALID_INPUT",
		})

	}

	result := action.DeviceRegistrationDB.Create(input)
	return res.Respond(result)
}

type FirebasePushRequest struct {
	UserId  string
	Content string
}

func PushFirebaseMessageToUser(req sdk.APIRequest, res sdk.APIResponder) error {
	var input FirebasePushRequest
	if err := req.GetContent(&input); err != nil {
		return res.Respond(&sdk.APIResponse{
			Status:    sdk.APIStatus.Invalid,
			ErrorCode: "FAIL_TO_PARSE_JSON",
			Message:   "Fail to parse JSON, please check the format. " + err.Error(),
		})
	}

	if len(input.UserId) == 0 || len(input.Content) == 0 {
		validateMsg := bytes.Buffer{}
		validateMsg.WriteString("require [UserId] [RegistrationId] [Content]")

		return res.Respond(&sdk.APIResponse{
			Status:    sdk.APIStatus.Invalid,
			Message:   validateMsg.String(),
			ErrorCode: "INVALID_INPUT",
		})
	}

	fbClient, err := action.GetDefaultFirebaseClient()
	if err != nil {
		return res.Respond(&sdk.APIResponse{
			Status:    sdk.APIStatus.Invalid,
			Message:   err.Error(),
			ErrorCode: "FAIL_CALL_FIREBASE",
		})
	}

	//
	filter := action.DeviceRegistration{
		UserId: input.UserId,
	}
	queryRs := action.DeviceRegistrationDB.Query(filter, 0, 100, true)
	if queryRs.Status != sdk.APIStatus.Ok {
		return res.Respond(queryRs)
	}

	devices := queryRs.Data.([]*action.DeviceRegistration)
	//

	var allResp []*action.FirebaseResponse

	for _, d := range devices {
		msg := []action.FirebaseMessages{
			fbClient.WithToken(d.RegistrationId),
			fbClient.WithData(map[string]interface{}{
				"message": input.Content,
			}),
		}

		fbResp, err := fbClient.Send(msg...)
		if err != nil {
			return res.Respond(&sdk.APIResponse{
				Status:    sdk.APIStatus.Invalid,
				Message:   err.Error(),
				ErrorCode: "FAIL_CALL_FIREBASE",
			})
		}

		if fbResp == nil {
			return res.Respond(&sdk.APIResponse{
				Status:    sdk.APIStatus.Invalid,
				ErrorCode: "FAIL_CALL_FIREBASE",
			})
		}

		allResp = append(allResp, fbResp)
	}

	return res.Respond(&sdk.APIResponse{
		Status:  sdk.APIStatus.Ok,
		Message: "push ok",
		Data:    allResp,
	})
}

func actionSendToFirebase(userId string, content string) error {

	fbClient, err := action.GetDefaultFirebaseClient()
	if err != nil {
		return err
	}

	filter := action.DeviceRegistration{
		UserId: userId,
	}
	queryRs := action.DeviceRegistrationDB.Query(filter, 0, 100, true)
	if queryRs.Status != sdk.APIStatus.Ok {
		return errors.New("Error: " + queryRs.Message)
	}

	// query devices
	devices := queryRs.Data.([]*action.DeviceRegistration)

	for _, d := range devices {
		msg := []action.FirebaseMessages{
			fbClient.WithToken(d.RegistrationId),
			fbClient.WithData(map[string]interface{}{
				"message": content,
			}),
		}

		_, err := fbClient.Send(msg...)
		if err != nil {
			return err
		}
	}

	return nil
}
