package api

import(
	"regexp"
	"progetto.wasa/service/database"
	"progetto.wasa/service/api/photoUtils"
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

func (group *Group) GroupFromDatabase(groupInDb database.Group) error {
	group.GroupId = groupInDb.GroupId
	group.Username = groupInDb.Username
	photo, err := photoUtils.ImageToBase64(photoUtils.GetProfilePhotoPath(group.GroupId))
	group.GroupPhoto = photo
	if err != nil {
		return err
	}
	return nil
}

func (group *Group) GroupToDatabase() database.Group {
	return database.Group{
		GroupId:   group.GroupId,
		Username: group.Username,
	}
}