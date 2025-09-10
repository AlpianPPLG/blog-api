package models

import (
	"time"

	"gorm.io/gorm"
)

type Post struct {
	ID        string    `json:"id" gorm:"primaryKey;type:varchar(36)"`
	Title     string    `json:"title" gorm:"type:varchar(255);not null"`
	Content   string    `json:"content" gorm:"type:text;not null"`
	Tags      []string  `json:"tags" gorm:"type:json"`
	AuthorID  string    `json:"author_id" gorm:"type:varchar(36);not null"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationships
	Author   User      `json:"author,omitempty" gorm:"foreignKey:AuthorID"`
	Comments []Comment `json:"comments,omitempty" gorm:"foreignKey:PostID"`
	Likes    []Like    `json:"likes,omitempty" gorm:"foreignKey:PostID"`
}

type PostCreateRequest struct {
	Title   string   `json:"title" binding:"required,min=1,max=255"`
	Content string   `json:"content" binding:"required,min=1"`
	Tags    []string `json:"tags"`
}

type PostUpdateRequest struct {
	Title   string   `json:"title" binding:"omitempty,min=1,max=255"`
	Content string   `json:"content" binding:"omitempty,min=1"`
	Tags    []string `json:"tags"`
}

type PostResponse struct {
	ID        string         `json:"id"`
	Title     string         `json:"title"`
	Content   string         `json:"content"`
	Tags      []string       `json:"tags"`
	AuthorID  string         `json:"author_id"`
	Author    UserResponse   `json:"author,omitempty"`
	Comments  []CommentResponse `json:"comments,omitempty"`
	Likes     []LikeResponse `json:"likes,omitempty"`
	LikesCount int           `json:"likes_count"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}
