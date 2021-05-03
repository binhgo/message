package main

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/binhgo/go-sdk/sdk"

	"github.com/binhgo/message/action"
	"github.com/binhgo/message/api"
	"github.com/binhgo/message/config"
	"github.com/binhgo/message/model"
)

var info = &apiInfo{}
var app *sdk.App

func onAllDBConnected() {
	model.MessageQueueDB.SetTopicConsumer(app.GetHostname(), action.QueueConsume)
	model.MessageQueueDB.StartConsume()

	// TEST()
}

func main() {

	info.StartTime = time.Now()
	info.Env = os.Getenv("env")
	info.Version = os.Getenv("version")

	// check sig
	signals = make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	go checkSig()
	// end

	// setup new app
	app = sdk.NewApp("Message Service")

	// get config & env
	configMap, err := app.GetConfigFromEnv()
	if err != nil {
		panic(err)
	}

	// init main DB
	var db = app.SetupDBClient(sdk.DBConfiguration{
		Address:  strings.Split(configMap["dbHost"], ","),
		Username: configMap["dbUser"],
		Password: configMap["dbPassword"],
		AuthDB:   config.Config.MainAuthDB,
	})
	db.OnConnected(onDBConnected)

	// init queue DB
	var queueHost = configMap["queueDbHost"]
	var queueUsername = configMap["queueUser"]
	var queuePassword = configMap["queuePassword"]
	if queueHost == "" {
		queueHost = configMap["dbHost"]
	}
	if queueUsername == "" {
		queueUsername = configMap["dbUser"]
	}
	if queuePassword == "" {
		queuePassword = configMap["dbPassword"]
	}
	db = app.SetupDBClient(sdk.DBConfiguration{
		Address:  strings.Split(queueHost, ","),
		Username: queueUsername,
		Password: queuePassword,
		AuthDB:   config.Config.QueueAuthDB,
	})
	db.OnConnected(onQueueDBConnected)

	// init log db
	var logHost = configMap["logDbHost"]
	var logUsername = configMap["logUser"]
	var logPassword = configMap["logPassword"]
	if logHost == "" {
		logHost = configMap["dbHost"]
	}
	if logUsername == "" {
		logUsername = configMap["dbUser"]
	}
	if logPassword == "" {
		logPassword = configMap["dbPassword"]
	}
	var db3 = app.SetupDBClient(sdk.DBConfiguration{
		Address:  strings.Split(logHost, ","),
		Username: logUsername,
		Password: logPassword,
		AuthDB:   config.Config.LogAuthDB,
	})
	db3.OnConnected(onLogDBConnected)

	// setup API Server
	protocol := os.Getenv("protocol")
	if protocol == "" {
		protocol = "THRIFT"
	}

	var server, _ = app.SetupAPIServer(protocol)
	server.SetHandler(sdk.APIMethod.GET, "/message/v1/api-info", healthCheck)
	// room
	server.SetHandler(sdk.APIMethod.POST, "/message/v1/room", api.CreateRoom)
	server.SetHandler(sdk.APIMethod.PUT, "/message/v1/room", api.UpdateRoomInfoPUT)
	server.SetHandler(sdk.APIMethod.PUT, "/message/v1/room/pin", api.PinMessage)
	server.SetHandler(sdk.APIMethod.GET, "/message/v1/room/pin", api.GetPinMessage)
	server.SetHandler(sdk.APIMethod.PUT, "/message/v1/room/add-user", api.AddUserToRoomPUT)
	server.SetHandler(sdk.APIMethod.PUT, "/message/v1/room/remove-user", api.RemoveUserToRoomPUT)
	server.SetHandler(sdk.APIMethod.GET, "/message/v1/room", api.GetRoom)
	server.SetHandler(sdk.APIMethod.GET, "/message/v1/room/info", api.GetRoomInfo)

	server.SetHandler(sdk.APIMethod.POST, "/message/v1/user/login", api.CreateUser)
	server.SetHandler(sdk.APIMethod.GET, "/message/v1/user/search", api.SearchUser)
	server.SetHandler(sdk.APIMethod.GET, "/message/v1/user/key", api.RequestUserPublicKey)

	// firebase
	server.SetHandler(sdk.APIMethod.POST, "/message/v1/firebase/device", api.RegisterNewFirebaseDevice)
	server.SetHandler(sdk.APIMethod.POST, "/message/v1/firebase/push", api.PushFirebaseMessageToUser)

	// message room
	server.SetHandler(sdk.APIMethod.GET, "/message/v1/message-room", api.GetMsgRoom)

	// signaling message
	server.SetHandler(sdk.APIMethod.POST, "/message/v1/message", api.CreateMessage)

	// expose
	server.Expose(80)

	// setup websocket server
	wss := app.SetupWSServer("message")
	wss.Timeout = config.Config.PingInterval + 10
	wss.Expose(8080)
	model.SetWebSocket(wss)
	wsRoute := wss.NewRoute("/")
	wsRoute.SetPayloadSize(4096 * 6)
	wsRoute.OnConnected = api.OnWSConnected
	wsRoute.OnMessage = api.OnWSMessage
	wsRoute.OnClose = api.OnWSClose

	// DB
	app.OnAllDBConnected(onAllDBConnected)

	// worker
	worker := app.SetupWorker()
	worker.SetTask(action.TryToCleanConnection)
	worker.SetDelay(100)
	worker.SetRepeatPeriod(config.Config.PingInterval)

	app.Launch()
}

type apiInfo struct {
	StartTime time.Time `json:"startTime"`
	Version   string    `json:"version"`
	Env       string    `json:"env"`
}

func healthCheck(req sdk.APIRequest, res sdk.APIResponder) error {
	return res.Respond(&sdk.APIResponse{
		Status:  sdk.APIStatus.Ok,
		Data:    []*apiInfo{info},
		Message: "Message Service",
	})
}

func onDBConnected(s *sdk.DBSession) error {

	fmt.Println("DB connected: " + config.Config.MainDBName)

	model.InitMessageDB(s, config.Config.MainDBName)
	model.InitUserConDB(s, config.Config.MainDBName)
	model.InitRoomDB(s, config.Config.MainDBName)
	model.InitMessageRoomDB(s, config.Config.MainDBName)
	model.InitUserDB(s, config.Config.MainDBName)
	model.InitKilledPod(s, config.Config.MainDBName)

	// firebase
	action.InitDeviceRegistrationDB(s, config.Config.MainDBName)

	return nil
}

func onQueueDBConnected(s *sdk.DBSession) error {

	fmt.Println("Queue connected: " + config.Config.QueueDBName)

	model.InitMessageQueueDB(s, config.Config.QueueDBName)
	action.InitFirebaseMessageQueueDB(s, config.Config.QueueDBName)

	return nil
}

func onLogDBConnected(s *sdk.DBSession) error {

	fmt.Println("Log connected: " + config.Config.LogDBName)

	model.InitMessageLogDB(s, config.Config.LogDBName)

	return nil
}
