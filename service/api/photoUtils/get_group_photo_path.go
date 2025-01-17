package photoUtils
import (
	"fmt"
)

func GetGroupPhotoPath(groupId int) string {
	return fmt.Sprintf("./storage/groups/%d/group_photo.jpg", groupId)
}