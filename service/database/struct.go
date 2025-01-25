package database

import "time"

type Message struct {
	MessageId int `json:"messageId"`
	DateTime time.Time `json:"dateTime"`
	Text string `json:"text"`
	Status bool `json:"status"`
	ConvId int `json:"convId"`
	SenderId int `json:"senderId"`
	Photo string `json:"photo"`
}

type User struct {
	UserId int `json:"userId"`
	Username string `json:"username"`
}

type Group struct {
	GroupId int `json:"groupId"`
	Username string `json:"username"`
}

type Conversation struct {
	ConvId int `json:"convId"`
	GroupId int `json:"groupId"`
	LastMessageId int `json:"lastMessageId"`
}