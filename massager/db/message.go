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
	tx := s.DB.Begin()

	if err := tx.Model(&models.Message{}).
		Where("receiver_number = ? AND sender_number = ?", userID1, userID2).
		Update("read", true).Error; err != nil {
		tx.Rollback()
		return nil, err
	}
	if err := tx.Where("(sender_number = ? AND receiver_number = ?) OR (sender_number = ? AND receiver_number = ?)",
		userID1, userID2, userID2, userID1).
		Order("created_at DESC").
		Find(&massage).Error; err != nil {
		tx.Rollback()
		return nil, err
	}
	tx.Commit()

	return massage, nil
}
