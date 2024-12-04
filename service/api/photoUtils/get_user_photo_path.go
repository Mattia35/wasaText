package photoUtils
import (
	"fmt"
)

func GetUserPhotoPath(userId int) string {
	return fmt.Sprintf("./storage/%d/user_photo.jpg", userId)
}