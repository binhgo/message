package client

import (
	"encoding/json"
	"os"
	"time"

	"gitlab.ghn.vn/internal-tools/message/config"

	"gitlab.ghn.vn/common-projects/go-sdk/sdk"
)

var ssoClient *sdk.RestClient

func getSSOClient() *sdk.RestClient {
	if ssoClient == nil {
		ssoClient = sdk.NewRESTClient(
			config.Config.SSO.Url,
			"sso_verify",
			10*time.Second,
			2,
			2*time.Second,
		)
		ssoClient.SetDebug(os.Getenv("env") != "prd")
	}
	return ssoClient
}

type SSOInfo struct {
	UserID   int    `json:"UserID"`
	Fullname string `json:"Fullname"`
	Phone    string `json:"Phone"`
}

type ssoResult struct {
	Code    int     `json:"code"`
	Message string  `json:"message"`
	Data    SSOInfo `json:"data"`
}

func VerifySSOToken(token string) *sdk.APIResponse {
	client := getSSOClient()

	// setup body
	params := make(map[string]string)
	params["token"] = token

	// setup headers
	headers := make(map[string]string)
	headers["Content-Type"] = "application/json"
	headers["Authorization"] = config.Config.SSO.Auth // config.Config.Key["sso-auth"]

	// call api
	result, err := client.MakeHTTPRequest(sdk.HTTPMethods.Get, headers, params, nil, "")
	if err != nil {
		return &sdk.APIResponse{
			Status:  sdk.APIStatus.Error,
			Message: "Endpoint error: " + err.Error(),
		}
	}

	// parse data
	var r = ssoResult{}
	err = json.Unmarshal(result.Content, &r)
	if err != nil || r.Code != 200 {
		message := r.Message
		if err != nil {
			message = err.Error()
		}
		return &sdk.APIResponse{
			Status:  sdk.APIStatus.Error,
			Message: "Parse result error: " + message,
		}
	}

	return &sdk.APIResponse{
		Status:  sdk.APIStatus.Ok,
		Message: "Get SSO successfully.",
		Data:    []SSOInfo{r.Data},
	}
}
