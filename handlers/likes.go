package handlers

import (
	"net/http"

	"blog-api/config"
	"blog-api/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// LikePost handles liking a post
func LikePost(c *gin.Context) {
	postID := c.Param("id")

	// Get user from context
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	userModel := user.(models.User)

	// Check if post exists
	var post models.Post
	if err := config.DB.First(&post, "id = ?", postID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	// Check if user already liked this post
	var existingLike models.Like
	if err := config.DB.Where("post_id = ? AND user_id = ?", postID, userModel.ID).First(&existingLike).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "You have already liked this post"})
		return
	}

	// Create like
	like := models.Like{
		ID:     uuid.New().String(),
		PostID: postID,
		UserID: userModel.ID,
	}

	if err := config.DB.Create(&like).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to like post"})
		return
	}

	// Load user information
	config.DB.Preload("User").First(&like, like.ID)

	likeResponse := models.LikeResponse{
		ID:        like.ID,
		PostID:    like.PostID,
		UserID:    like.UserID,
		User: models.UserResponse{
			ID:        like.User.ID,
			Username:  like.User.Username,
			Email:     like.User.Email,
			CreatedAt: like.User.CreatedAt,
		},
		CreatedAt: like.CreatedAt,
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Post liked successfully",
		"like":    likeResponse,
	})
}

// UnlikePost handles unliking a post
func UnlikePost(c *gin.Context) {
	postID := c.Param("id")

	// Get user from context
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	userModel := user.(models.User)

	// Check if post exists
	var post models.Post
	if err := config.DB.First(&post, "id = ?", postID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	// Find existing like
	var like models.Like
	if err := config.DB.Where("post_id = ? AND user_id = ?", postID, userModel.ID).First(&like).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "You have not liked this post"})
		return
	}

	// Delete like
	if err := config.DB.Delete(&like).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unlike post"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Post unliked successfully",
	})
}

// GetPostLikes handles getting all likes for a post
func GetPostLikes(c *gin.Context) {
	postID := c.Param("id")

	// Check if post exists
	var post models.Post
	if err := config.DB.First(&post, "id = ?", postID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	// Get likes with user information
	var likes []models.Like
	if err := config.DB.Preload("User").Where("post_id = ?", postID).Find(&likes).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch likes"})
		return
	}

	// Convert to response format
	var likesResponse []models.LikeResponse
	for _, like := range likes {
		likeResponse := models.LikeResponse{
			ID:        like.ID,
			PostID:    like.PostID,
			UserID:    like.UserID,
			User: models.UserResponse{
				ID:        like.User.ID,
				Username:  like.User.Username,
				Email:     like.User.Email,
				CreatedAt: like.User.CreatedAt,
			},
			CreatedAt: like.CreatedAt,
		}
		likesResponse = append(likesResponse, likeResponse)
	}

	c.JSON(http.StatusOK, gin.H{
		"likes": likesResponse,
		"count": len(likesResponse),
	})
}

// CheckUserLike handles checking if a user has liked a post
func CheckUserLike(c *gin.Context) {
	postID := c.Param("id")

	// Get user from context
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	userModel := user.(models.User)

	// Check if user liked this post
	var like models.Like
	if err := config.DB.Where("post_id = ? AND user_id = ?", postID, userModel.ID).First(&like).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"liked": false,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"liked": true,
		"like_id": like.ID,
		"liked_at": like.CreatedAt,
	})
}
