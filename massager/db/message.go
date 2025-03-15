package db

import (
	models "chach/massager/db/model"
)

func (s *Storege) Send(massages *models.Message) error {

	err := s.DB.Create(&massages).Error
	if err != nil {
		return err
	}
	return nil
}
