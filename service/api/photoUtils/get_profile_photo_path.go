package photoUtils

import (
	"fmt"
)

func GetProfilePhotoPath(userId int) string {
	return fmt.Sprintf("./storage/%d/user_propic_250x250.jpg", userId)
}