package models

import (
	"time"

	"gorm.io/gorm"
)

type Comment struct {
	ID              string    `json:"id" gorm:"primaryKey;type:varchar(36)"`
	PostID          string    `json:"post_id" gorm:"type:varchar(36);not null"`
	AuthorID        string    `json:"author_id" gorm:"type:varchar(36);not null"`
	ParentCommentID *string   `json:"parent_comment_id" gorm:"type:varchar(36);null"`
	Content         string    `json:"content" gorm:"type:text;not null"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationships
	Post            Post      `json:"post,omitempty" gorm:"foreignKey:PostID"`
	Author          User      `json:"author,omitempty" gorm:"foreignKey:AuthorID"`
	ParentComment   *Comment  `json:"parent_comment,omitempty" gorm:"foreignKey:ParentCommentID"`
	Replies         []Comment `json:"replies,omitempty" gorm:"foreignKey:ParentCommentID"`
}

type CommentCreateRequest struct {
	Content string `json:"content" binding:"required,min=1"`
}

type CommentReplyRequest struct {
	Content string `json:"content" binding:"required,min=1"`
}

type CommentUpdateRequest struct {
	Content string `json:"content" binding:"required,min=1"`
}

type CommentResponse struct {
	ID              string            `json:"id"`
	PostID          string            `json:"post_id"`
	AuthorID        string            `json:"author_id"`
	ParentCommentID *string           `json:"parent_comment_id"`
	Content         string            `json:"content"`
	Author          UserResponse      `json:"author,omitempty"`
	Replies         []CommentResponse `json:"replies,omitempty"`
	CreatedAt       time.Time         `json:"created_at"`
	UpdatedAt       time.Time         `json:"updated_at"`
}
