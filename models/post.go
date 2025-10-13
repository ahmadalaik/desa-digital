package models

import "time"

type Post struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	Image      string    `json:"image"`
	Title      string    `json:"title"`
	Slug       string    `json:"slug" gorm:"unique"`
	Content    string    `json:"content"`
	CategoryId string    `json:"category_id"`
	Category   string    `json:"category" gorm:"foreignKey:CategoryId"`
	UserId     uint      `json:"user_id"`
	User       User      `json:"user" gorm:"foreignKey:UserId"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
