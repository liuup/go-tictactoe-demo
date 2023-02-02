package models

import (
	"errors"
	"main/utils/token"
)

type User struct {
	ID       uint
	Name     string
	Password string
	Age      int
	Tel      string
	// updateAt *time.Time
}

func GetUserByID(uid uint) (User, error) {
	u := User{}

	if err := DB.First(&u, uid).Error; err != nil {
		return u, errors.New("User not found!")
	}

	u.PrepareGive()

	return u, nil
}

func (u *User) PrepareGive() {
	u.Password = ""
}

// 登录检查
func LoginCheck(name string, password string) (string, error) {
	u := User{}

	err := DB.Model(User{}).Where("name = ?", name).Take(&u).Error
	if err != nil {
		return "", err
	}

	// 校验密码
	if u.Password != password {
		return "", errors.New("password incorrect")
	}

	token, err := token.GenerateToken(u.ID)
	if err != nil {
		return "", err
	}

	return token, nil
}

// 注册保存用户
func (u *User) SaveUser() (*User, error) {
	err := DB.Select("Name", "Password").Create(&u).Error
	// err := DB.Create(&User{Name: u.Name, Password: u.Password}).Error
	// err := DB.Create(&u).Error
	if err != nil {
		return &User{}, err
	}

	return u, nil
}

// 注册数据的预处理
// func (u *User) BeforeSave() error {
// 	hashedPassword, err :=
// }
