package db

import (
	"chach/massager/db/model"
)

func (s *Storege) CreatUser(user *model.User) error {

	err := s.DB.Create(user).Error
	if err != nil {
		return err
	}
	return nil
}
