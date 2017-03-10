package services

import (
	"github.com/krufyliu/dkvgo/dkvgo-admin/models"
)

type userService struct {} 

func (this *userService) GetUserById(userId int) (*models.User, error) {
	return nil, nil
}

func (this *userService) GetUserByEmail(email string) (*models.User, error) {
	user := &models.User{}
	user.Email = email
	err := o.Read(user, "Email")
	if err != nil {
		return nil, err
	}
	return user, nil
}
