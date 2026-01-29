package main

import "gorm.io/gorm"

type PostRepository struct {
	DB *gorm.DB
}

func (r *PostRepository) GetAll() ([]PostDB, error) {
	var posts []PostDB
	err := r.DB.Find(&posts).Error
	if err != nil {
		return nil, &ErrorSchema{Code: 500, Text: "DB error: " + err.Error()}
	}
	if len(posts) == 0 {
		return nil, ErrPostNotFound
	}
	return posts, nil
}
