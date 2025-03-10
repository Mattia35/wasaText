package structions



type Comment struct {
	MessageId int `json:"messageId"`
	Content string `json:"content"`
	SenderId int `json:"senderId"`
	ConvId int `json:"convId"`
}