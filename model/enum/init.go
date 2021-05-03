package enum

type TopicEnumValue string

type TopicEnum struct {
	NONE              TopicEnumValue
	CONNECTION        TopicEnumValue
	PING              TopicEnumValue
	AUTHORIZATION     TopicEnumValue
	MESSAGE           TopicEnumValue
	IMAGE             TopicEnumValue
	FILE              TopicEnumValue
	ROOM_NOTIFY       TopicEnumValue
	VIDEO_OFFER       TopicEnumValue
	VIDEO_ANSWER      TopicEnumValue
	NEW_ICE_CANDIDATE TopicEnumValue
	HANG_UP			  TopicEnumValue
}

var Topic = &TopicEnum{
	NONE:              TopicEnumValue("NONE"), // don't use. Just used for upsert & get new id.
	CONNECTION:        TopicEnumValue("CONNECTION"),
	PING:              TopicEnumValue("PING"),
	AUTHORIZATION:     TopicEnumValue("AUTHORIZATION_V2"),
	MESSAGE:           TopicEnumValue("MESSAGE"),
	IMAGE:             TopicEnumValue("IMAGE"),
	FILE:              TopicEnumValue("FILE"),
	ROOM_NOTIFY:       TopicEnumValue("ROOM_NOTIFY"),
	VIDEO_OFFER:       TopicEnumValue("VIDEO_OFFER"),
	VIDEO_ANSWER:      TopicEnumValue("VIDEO_ANSWER"),
	NEW_ICE_CANDIDATE: TopicEnumValue("NEW_ICE_CANDIDATE"),
	HANG_UP: 		   TopicEnumValue("HANG_UP"),
}

type RoomTypeEnumValue string
type RoomTypeEnum struct {
	NONE RoomTypeEnumValue
	ONE  RoomTypeEnumValue
	MANY RoomTypeEnumValue
}

var RoomType = &RoomTypeEnum{
	NONE: RoomTypeEnumValue("NONE"),
	MANY: RoomTypeEnumValue("MANY"),
	ONE:  RoomTypeEnumValue("ONE"),
}

type ConStatusEnumValue string

type ConStatusEnum struct {
	ACTIVE ConStatusEnumValue
	CLOSED ConStatusEnumValue
}

var ConStatus = &ConStatusEnum{
	ACTIVE: ConStatusEnumValue("ACTIVE"),
	CLOSED: ConStatusEnumValue("CLOSED"),
}

type MessageTypeEnumValue string
type MessageTypeEnum struct {
	RECEIVE MessageTypeEnumValue
	PUSH    MessageTypeEnumValue
}

var MessageType = &MessageTypeEnum{
	RECEIVE: MessageTypeEnumValue("RECEIVE"),
	PUSH:    MessageTypeEnumValue("PUSH"),
}

type MessageCategoryEnum string

const (
	PUSH_MESSAGE = MessageCategoryEnum("PUSH_MESSAGE")
)

type MessageStatusEnumValue string
type MessageStatusEnum struct {
	DELIVERING MessageStatusEnumValue
	DELIVERED  MessageStatusEnumValue
}

var MessageStatus = &MessageStatusEnum{
	DELIVERING: MessageStatusEnumValue("DELIVERING"),
	DELIVERED:  MessageStatusEnumValue("DELIVERED"),
}

type DeliverTypeEnumValue string
type deliverTypeEnum struct {
	DIRECT DeliverTypeEnumValue
	PULL   DeliverTypeEnumValue
}

var DeliverType = &deliverTypeEnum{
	DIRECT: DeliverTypeEnumValue("DIRECT"),
	PULL:   DeliverTypeEnumValue("PULL"),
}
