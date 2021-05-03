package client

import (
	"encoding/json"
	"os"
	"time"

	"github.com/binhgo/go-sdk/sdk"

	"github.com/binhgo/message/config"
)

var ssov2Client *sdk.RestClient

func getSSOv2Client() *sdk.RestClient {
	if ssoClient == nil {
		ssoClient = sdk.NewRESTClient(
			config.Config.SSOv2.Url,
			"sso_v2_verify",
			10*time.Second,
			2,
			2*time.Second,
		)
		ssoClient.SetDebug(os.Getenv("env") != "prd")
	}
	return ssoClient
}

type VerifyTokenSSOv2Data struct {
	UserID int    `json:"user_id"`
	Name   string `json:"name"`
	Phone  string `json:"phone"`
}

type VerifyTokenSSOv2Result struct {
	Code    int                  `json:"code"`
	Message string               `json:"message"`
	Data    VerifyTokenSSOv2Data `json:"data"`
}

func VerifySSOv2Token(token, userAgent, IP string) *sdk.APIResponse {
	client := getSSOv2Client()

	// setup body
	body := make(map[string]string)
	body["access_token"] = token
	body["app_key"] = config.Config.SSOv2.AppKey
	body["app_secret"] = config.Config.SSOv2.AppSecret
	body["user_agent"] = userAgent
	body["remote_ip"] = IP

	// setup headers
	headers := make(map[string]string)
	headers["Content-Type"] = "application/json"

	// call api
	result, err := client.MakeHTTPRequest(sdk.HTTPMethods.Post, headers, nil, body, "/verify-access-token")
	if err != nil {
		return &sdk.APIResponse{
			Status:  sdk.APIStatus.Error,
			Message: "Endpoint error: " + err.Error(),
		}
	}

	// parse data
	var r = VerifyTokenSSOv2Result{}
	err = json.Unmarshal(result.Content, &r)
	if err != nil || r.Code != 200 {
		message := r.Message
		if err != nil {
			message = err.Error()
		}
		return &sdk.APIResponse{
			Status:  sdk.APIStatus.Error,
			Message: "Verify Token error: " + message,
		}
	}

	return &sdk.APIResponse{
		Status:  sdk.APIStatus.Ok,
		Message: "Verify Token successfully.",
		Data:    []VerifyTokenSSOv2Data{r.Data},
	}
}

type GenTokenSSOv2Data struct {
	AccessToken string `json:"access_token"`
}

type GenTokenSSOv2Result struct {
	Code    int               `json:"code"`
	Message string            `json:"message"`
	Data    GenTokenSSOv2Data `json:"data"`
}

// GenTokenSSOv2 ...
func GenTokenSSOv2(authorKey, userAgent, IP string) *sdk.APIResponse {
	client := getSSOv2Client()

	headers := make(map[string]string)
	headers["Content-Type"] = "application/json"

	// setup body
	body := make(map[string]string)
	body["authorization_code"] = authorKey
	body["app_key"] = config.Config.SSOv2.AppKey
	body["app_secret"] = config.Config.SSOv2.AppSecret
	body["user_agent"] = userAgent
	body["remote_ip"] = IP

	// call api
	result, err := client.MakeHTTPRequest(sdk.HTTPMethods.Post, headers, nil, body, "/gen-access-token")
	if err != nil {
		return &sdk.APIResponse{
			Status:  sdk.APIStatus.Error,
			Message: "Endpoint error: " + err.Error(),
		}
	}

	// parse data
	var r = GenTokenSSOv2Result{}
	err = json.Unmarshal(result.Content, &r)
	if err != nil || r.Code != 200 {
		message := r.Message
		if err != nil {
			message = err.Error()
		}
		return &sdk.APIResponse{
			Status:  sdk.APIStatus.Error,
			Message: "Gen Token error: " + message,
		}
	}

	return &sdk.APIResponse{
		Status:  sdk.APIStatus.Ok,
		Message: "Gen token successfully.",
		Data:    []GenTokenSSOv2Data{r.Data},
	}
}
