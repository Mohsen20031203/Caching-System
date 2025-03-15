package db

import (
	models "chach/massager/db/model"
)

func (s *Storege) Getmassages(id int) ([]models.Message, error) {
	var massages models.Message

	err := s.DB.First(id, &massages).Error
	if err != nil {
		return nil, err
	}
	return nil, nil
}
