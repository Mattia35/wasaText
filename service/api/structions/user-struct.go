package structions

import(
	"regexp"
)

type User struct {
	UserId int `json:"userId"`
	Username string `json:"username"`
	UserPhoto string `json:"userPhoto"`
} 

func (user *User) IsValid() bool {
	username := user.Username
	validUser := regexp.MustCompile(`^.*$`)
	return validUser.MatchString(username)
}
