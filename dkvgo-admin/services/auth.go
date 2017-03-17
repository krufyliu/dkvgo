package services

import "github.com/krufyliu/dkvgo/dkvgo-admin/models"

type authService struct {
	user *models.User
}

func (this *authService) IsLogined() bool {
	return this.user != nil
}

func (this *authService) GetAuthUser() *models.User {
	return this.user
}

func (this *authService) LoginUser(user *models.User) {
	this.user = user
}