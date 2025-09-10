package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID           string    `json:"id" gorm:"primaryKey;type:varchar(36)"`
	Username     string    `json:"username" gorm:"uniqueIndex;type:varchar(50);not null"`
	Email        string    `json:"email" gorm:"uniqueIndex;type:varchar(100);not null"`
	PasswordHash string    `json:"-" gorm:"type:varchar(255);not null"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationships
	Posts    []Post    `json:"posts,omitempty" gorm:"foreignKey:AuthorID"`
	Comments []Comment `json:"comments,omitempty" gorm:"foreignKey:AuthorID"`
	Likes    []Like    `json:"likes,omitempty" gorm:"foreignKey:UserID"`
}

type UserRegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type UserLoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type UserResponse struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}
