package structions



type Comment struct {
	CommentId int `json:"commentId"`
	MessageId int `json:"messageId"`
	Content string `json:"content"`
	SenderId int `json:"senderId"`
	ConvId int `json:"convId"`
}