package model

import (
	"github.com/binhgo/go-sdk/sdk"

	"github.com/binhgo/message/model/enum"
)

type MessageQueueItem struct {
	Category enum.MessageCategoryEnum
	/**
		-Data depend on -Category
	    PUSH_MESSAGE => -Data = message id
	*/
	Data              string
	ConnectionLocalID int
}

var MessageQueueDB = sdk.DBSortedQueue{
	ColName: "message_queue",
}

func InitMessageQueueDB(dbSession *sdk.DBSession, dbName string) {
	MessageQueueDB.Init(dbSession, dbName)
}
