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

func (s *Storege) Read(id string) error {
	if err := s.DB.Model(&models.Message{}).Where("id = ?", id).Update("read", true).Error; err != nil {

		return err
	}
	return nil
}
