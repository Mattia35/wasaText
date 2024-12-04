package api

import(
	"regexp"
	"progetto.wasa/service/database"
	"progetto.wasa/service/api/photoUtils"
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

func (user *User) FromDatabase(userInDb database.User) error {
	user.UserId = userInDb.UserId
	user.Username = userInDb.Username
	photo, err := photoUtils.ImageToBase64(photoUtils.GetProfilePhotoPath(user.UserId))
	user.UserPhoto = photo
	if err != nil {
		return err
	}
	return nil
}

func (u *User) ToDatabase() database.User {
	return database.User{
		UserId:   u.UserId,
		Username: u.Username,
	}
}