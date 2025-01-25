package structions

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