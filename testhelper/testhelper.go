package testhelper

import "github.com/google/uuid"

func GetTempFileLoc() string {
	return "/tmp/" + uuid.NewString()
}
