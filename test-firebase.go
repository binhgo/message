package main

import (
	"gitlab.ghn.vn/internal-tools/message/action"
)

func test_firebase_send_msg() {

	c, err := action.GetDefaultFirebaseClient()
	if err != nil {
		panic(err)
	}

	msg := []action.FirebaseMessages{
		c.WithToken("fcm-token"),
		c.WithNotification(&action.FirebaseNotification{
			Title: "testtitle",
			Body:  "testBody",
		}),
		c.WithData(map[string]interface{}{
			"message": "test---data",
		}),
	}

	res, err := c.Send(msg...)

	if err != nil {
		panic(err)
	}

	if res == nil {
		panic("No response from the server")
	}
}
