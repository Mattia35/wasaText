package structions

type Conversation struct {
	convId int `json:"convId"`
	lastMessage int `json:"lastMessage"`
	group bool `json:"group"`
}