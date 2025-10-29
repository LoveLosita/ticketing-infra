package dao

import (
	"ticketing-infra/rpc-server/inits"
	"ticketing-infra/rpc-server/user-service/model"
)

func InsertNewUserInfo(user model.User) (int, error) {
	result := inits.Db.Create(&user)
	if result.Error != nil {
		return 0, result.Error
	}
	return int(user.ID), nil
}

func IfUsernameExists(name string) (bool, error) {
	var nameUser model.User
	result := inits.Db.First(&nameUser, "username = ?", name)
	if result.Error != nil {
		if result.Error.Error() == "record not found" {
			return false, nil
		}
		return false, result.Error
	}
	return true, nil
}
