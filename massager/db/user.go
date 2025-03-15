package db

import (
	models "chach/massager/db/model"
)

func (s *Storege) CreatUser(user *models.User) error {

	err := s.DB.Create(user).Error
	if err != nil {
		return err
	}
	return nil
}

func (s *Storege) GetUser(id int64) (*models.User, error) {
	var user models.User

	err := s.DB.First(&user, id).Error
	if err != nil {
		return nil, err
	}

	return &user, err
}

func (s *Storege) GetUsers() ([]models.User, error) {
	var users []models.User
	err := s.DB.Order("id desc").Find(&users).Error
	return users, err
}
