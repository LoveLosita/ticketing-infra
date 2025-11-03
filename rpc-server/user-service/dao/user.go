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

func GetUserHashedPassword(username string) (string, error) {
	var user model.User
	result := inits.Db.First(&user, "username = ?", username)
	if result.Error != nil {
		if result.Error.Error() == "record not found" {
			return "", nil
		}
		return "", result.Error
	}
	return user.Password, nil
}

func GetUserIDByUsername(username string) (int, error) {
	var user model.User
	result := inits.Db.First(&user, "username = ?", username)
	if result.Error != nil {
		return 0, result.Error
	}
	return int(user.ID), nil
}

func ChangeUserPassword(userName string, newHashedPwd string) error {
	result := inits.Db.Model(&model.User{}).Where("username = ?", userName).Update("password", newHashedPwd)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func IfUserIDExists(userID int) (bool, error) {
	var user model.User
	result := inits.Db.First(&user, "id = ?", userID)
	if result.Error != nil {
		if result.Error.Error() == "record not found" {
			return false, nil
		}
		return false, result.Error
	}
	return true, nil
}

func SetUserRoleToAdmin(userID int) error {
	result := inits.Db.Model(&model.User{}).Where("id = ?", userID).Update("role", "admin")
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func GetUserRoleByID(userID int) (string, error) {
	var user model.User
	result := inits.Db.First(&user, "id = ?", userID)
	if result.Error != nil {
		return "", result.Error
	}
	return user.Role, nil
}
