package database

type Message struct {
	MessageId int `json:"messageId"`
	Text string `json:"text"`
	Status bool `json:"status"`
	ConvId int `json:"convId"`
	SenderId int `json:"senderId"`
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
	OtherUserId int `json:"otherUserId"`
	SenderId int `json:"senderId"`
	LastMessageId int `json:"lastMessageId"`
}