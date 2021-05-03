package main

import (
	"fmt"
	"os"
	"syscall"
	"time"

	"github.com/binhgo/go-sdk/sdk"

	"github.com/binhgo/message/model"
)

var signals chan os.Signal

func checkSig() {
	sig := <-signals
	fmt.Println(sig)
	if sig == syscall.SIGTERM {
		onSigTerm()
	}
}

func onSigTerm() {

	fmt.Println("////////////////////")
	fmt.Println("SIGTERM")
	fmt.Println("////////////////////")

	now := time.Now()

	killedPod := model.KilledPod{
		PodName:  app.GetHostname(),
		FailTime: &now,
		Reason:   "SIGTERM",
	}

	result := model.KilledPodDB.Create(killedPod)
	if result.Status != sdk.APIStatus.Ok {
		fmt.Println("////////////////////")
		fmt.Println(killedPod)
		fmt.Println("////////////////////")
	}

	time.Sleep(100 * time.Millisecond)
}
