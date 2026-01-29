package main

import "time"

type PostRequest struct {
	Title string  `json:"title" binding:"required"`
	Body  *string `json:"body"`
}

type Post struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Body  string `json:"body"`
}

type PostDB struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

func (PostDB) TableName() string {
	return "posts"
}
