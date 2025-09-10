package models

import (
	"time"
)

type Like struct {
	ID        string    `json:"id" gorm:"primaryKey;type:varchar(36)"`
	PostID    string    `json:"post_id" gorm:"type:varchar(36);not null"`
	UserID    string    `json:"user_id" gorm:"type:varchar(36);not null"`
	CreatedAt time.Time `json:"created_at"`

	// Relationships
	Post Post `json:"post,omitempty" gorm:"foreignKey:PostID"`
	User User `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

type LikeResponse struct {
	ID        string    `json:"id"`
	PostID    string    `json:"post_id"`
	UserID    string    `json:"user_id"`
	User      UserResponse `json:"user,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}
