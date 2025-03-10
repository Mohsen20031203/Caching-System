package db

import "chach/massager/db/model"

func (s *Storege) Getmassages(id int) ([]model.Message, error) {
	var massages model.Message

	err := s.DB.First(id, &massages).Error
	if err != nil {
		return nil, err
	}
	return nil, nil
}
