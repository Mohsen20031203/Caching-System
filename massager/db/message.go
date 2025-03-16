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

func (s *Storege) GetMessagesBetweenUsers(userID1, userID2 uint) ([]models.Message, error) {
	var massage []models.Message
	err := s.DB.Where("(sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?)",
		userID1, userID2, userID2, userID1).
		Order("created_at DESC").
		Find(&massage).Error

	if err != nil {
		return nil, err
	}
	return massage, nil
}
