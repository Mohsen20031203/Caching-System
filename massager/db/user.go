package db

import (
	models "chach/massager/db/model"
	"fmt"
	"log"
)

func (s *Storege) SignUp(user *models.User) error {

	if err := s.DB.Where("phone = ?", user.Phone).First(&user).Error; err == nil {

		return fmt.Errorf("User already exists")
	}

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

func (s *Storege) GetUsers() ([]models.User, error) {
	var users []models.User
	err := s.DB.Order("id desc").Find(&users).Error
	return users, err
}

func (s *Storege) DeleteUser(id uint) error {
	var user models.User
	if err := s.DB.Where("id = ?", id).First(&user).Error; err != nil {
		//ctx.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
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
