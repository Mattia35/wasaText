package structions

import(
	"regexp"
)



type Group struct {
	GroupId int `json:"groupId"`
	Username string `json:"username"`
	GroupPhoto string `json:"groupPhoto"`
}

func (group *Group) IsValid() bool {
	username := group.Username
	validGroup := regexp.MustCompile(`^.*$`)
	return validGroup.MatchString(username)
}
