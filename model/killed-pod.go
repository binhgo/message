package model

import (
	"time"

	"gitlab.ghn.vn/common-projects/go-sdk/sdk"
)

type KilledPod struct {
	PodName  string
	Reason   string
	FailTime *time.Time
}

var KilledPodDB = &sdk.DBModel2{
	ColName:        "killed_pod",
	TemplateObject: &KilledPod{},
}

func InitKilledPod(dbSession *sdk.DBSession, dbName string) {
	KilledPodDB.DBName = dbName
	KilledPodDB.Init(dbSession)
}
