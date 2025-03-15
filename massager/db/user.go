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
