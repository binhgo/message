package api

import (
	"fmt"

	"github.com/binhgo/go-sdk/sdk"
	"github.com/globalsign/mgo/bson"

	"github.com/binhgo/message/model"
)

func CreateUser(req sdk.APIRequest, res sdk.APIResponder) error {

	var input model.User

	if err := req.GetContent(&input); err != nil {
		return res.Respond(&sdk.APIResponse{
			Status:    sdk.APIStatus.Invalid,
			ErrorCode: "FAIL_TO_PARSE_JSON",
			Message:   "Fail to parse JSON, please check the format. " + err.Error(),
		})
	}

	if len(input.UserID) == 0 || len(input.Password) == 0 {
		return res.Respond(&sdk.APIResponse{
			Status:    sdk.APIStatus.Invalid,
			ErrorCode: "INVALID_INFORMATION",
			Message:   "INVALID: body input.",
		})
	}

	filter := &model.User{
		UserID: input.UserID,
	}

	queryRs := model.UserDB.QueryOne(filter)

	switch queryRs.Status {
	case sdk.APIStatus.Ok:
		userInfo := queryRs.Data.([]*model.User)[0]
		if input.Password != userInfo.Password {
			return res.Respond(&sdk.APIResponse{
				Status:    sdk.APIStatus.Forbidden,
				ErrorCode: "WRONG_PASSWORD",
				Message:   "Password doesn't match.",
			})
		}
		return res.Respond(queryRs)
	case sdk.APIStatus.NotFound:
		createRs := model.UserDB.Create(input)
		return res.Respond(createRs)
	default:
		return res.Respond(queryRs)
	}
}

func SearchUser(req sdk.APIRequest, res sdk.APIResponder) error {
	userId := req.GetParam("userId")

	if userId == "" || len(userId) == 0 {
		return res.Respond(&sdk.APIResponse{
			Status:    sdk.APIStatus.Invalid,
			ErrorCode: "INVALID_PARAMS",
			Message:   "Invalid [userId] params.",
		})
	}

	searchFilter := bson.M{}

	pat := fmt.Sprintf("%s.*", userId)
	searchFilter["user_id"] = bson.RegEx{
		Pattern: pat,
		Options: "i",
	}

	resp := model.UserDB.Query(searchFilter, 0, 0, true)

	return res.Respond(resp)
}

func RequestUserPublicKey(req sdk.APIRequest, res sdk.APIResponder) error {
	var input model.User

	if err := req.GetContent(&input); err != nil {
		return res.Respond(&sdk.APIResponse{
			Status:    sdk.APIStatus.Invalid,
			ErrorCode: "FAIL_TO_PARSE_JSON",
			Message:   "Fail to parse JSON, please check the format. " + err.Error(),
		})
	}

	if len(input.UserID) == 0 {
		return res.Respond(&sdk.APIResponse{
			Status:    sdk.APIStatus.Invalid,
			ErrorCode: "FAIL_TO_PARSE_JSON",
			Message:   "require [UserID]",
		})
	}

	filter := &model.User{
		UserID: input.UserID,
	}

	queryRs := model.UserDB.QueryOne(filter)
	if queryRs.Status != sdk.APIStatus.Ok {
		return res.Respond(&sdk.APIResponse{
			Status:    queryRs.Status,
			ErrorCode: "ERR_GET_PUB_KEY",
			Message:   queryRs.Message,
		})
	}

	user := queryRs.Data.([]*model.User)[0]
	if len(user.PubKey) == 0 {
		return res.Respond(&sdk.APIResponse{
			Status:    sdk.APIStatus.Invalid,
			ErrorCode: "ERR_GET_PUB_KEY",
			Message:   "pubKey empty",
		})
	}

	return res.Respond(&sdk.APIResponse{
		Status: sdk.APIStatus.Ok,
		Data:   []string{user.PubKey},
	})
}
