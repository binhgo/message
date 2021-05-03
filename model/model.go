package model

import (
	"encoding/json"

	"github.com/binhgo/message/model/enum"
)

type WSInputMessage struct {
	Topic     enum.TopicEnumValue `json:"topic,omitempty" bson:"topic,omitempty"`
	Content   *Content            `json:"content,omitempty" bson:"content,omitempty"`
	UserAgent string              `json:"userAgent,omitempty" bson:"userAgent,omitempty"`
}

type WSOutputMessage struct {
	Topic   enum.TopicEnumValue `json:"topic,omitempty" bson:"topic,omitempty"`
	Content *Content            `json:"content,omitempty" bson:"content,omitempty"`
}

func (outMsg *WSOutputMessage) String() string {
	bytes, err := json.Marshal(outMsg)
	if err != nil {
		return ""
	}
	return string(bytes)
}
