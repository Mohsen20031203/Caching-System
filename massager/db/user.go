package db

import (
	models "chach/massager/db/model"
	"fmt"
	"log"
)

func (s *Storege) SignUp(user *models.User) error {

	err := s.DB.Create(user).Error
	if err != nil {
		return err
	}
	return nil
}

func (s *Storege) GetUser(number string) (*models.User, error) {
	var user models.User

	err := s.DB.Where("phone = ?", number).First(&user).Error
	if err != nil {
		return nil, err
	}

	return &user, err
}

func (s *Storege) GetUsers(number string) ([]models.User, error) {
	var users []models.User

	// Find users who either sent messages to the specified number or received messages from the specified number
	err := s.DB.Joins("JOIN messages ON messages.sender_number = users.phone OR messages.receiver_number = users.phone").
		Where("messages.sender_number = ? OR messages.receiver_number = ?", number, number).
		Where("users.phone != ?", number).
		Distinct().
		Find(&users).Error

	if err != nil {
		return nil, err
	}

	return users, nil
}

func (s *Storege) DeleteUser(id uint) error {
	var user models.User
	if err := s.DB.Where("id = ?", id).First(&user).Error; err != nil {
		return err
	}

	if err := s.DB.Delete(&user).Error; err != nil {
		log.Println("Error deleting user:", err)
		return err
	}
	return nil
}

func (s *Storege) UpdateUser(user models.User) error {

	if err := s.DB.Save(&user).Error; err != nil {
		return fmt.Errorf("Failed to update user")
	}
	return nil
}
