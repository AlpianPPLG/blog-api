package handlers

import (
	"net/http"

	"blog-api/config"
	"blog-api/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// CreateComment handles creating a comment on a post
func CreateComment(c *gin.Context) {
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

	// Parse comment request
	var req models.CommentCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create comment
	comment := models.Comment{
		ID:       uuid.New().String(),
		PostID:   postID,
		AuthorID: userModel.ID,
		Content:  req.Content,
	}

	if err := config.DB.Create(&comment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create comment"})
		return
	}

	// Load author information
	config.DB.Preload("Author").First(&comment, comment.ID)

	commentResponse := models.CommentResponse{
		ID:              comment.ID,
		PostID:          comment.PostID,
		AuthorID:        comment.AuthorID,
		ParentCommentID: comment.ParentCommentID,
		Content:         comment.Content,
		Author: models.UserResponse{
			ID:        comment.Author.ID,
			Username:  comment.Author.Username,
			Email:     comment.Author.Email,
			CreatedAt: comment.Author.CreatedAt,
		},
		CreatedAt: comment.CreatedAt,
		UpdatedAt: comment.UpdatedAt,
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Comment created successfully",
		"comment": commentResponse,
	})
}

// ReplyToComment handles replying to a comment
func ReplyToComment(c *gin.Context) {
	commentID := c.Param("id")

	// Get user from context
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	userModel := user.(models.User)

	// Check if parent comment exists
	var parentComment models.Comment
	if err := config.DB.First(&parentComment, "id = ?", commentID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Parent comment not found"})
		return
	}

	// Parse reply request
	var req models.CommentReplyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create reply comment
	comment := models.Comment{
		ID:              uuid.New().String(),
		PostID:          parentComment.PostID,
		AuthorID:        userModel.ID,
		ParentCommentID: &commentID,
		Content:         req.Content,
	}

	if err := config.DB.Create(&comment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create reply"})
		return
	}

	// Load author information
	config.DB.Preload("Author").First(&comment, comment.ID)

	commentResponse := models.CommentResponse{
		ID:              comment.ID,
		PostID:          comment.PostID,
		AuthorID:        comment.AuthorID,
		ParentCommentID: comment.ParentCommentID,
		Content:         comment.Content,
		Author: models.UserResponse{
			ID:        comment.Author.ID,
			Username:  comment.Author.Username,
			Email:     comment.Author.Email,
			CreatedAt: comment.Author.CreatedAt,
		},
		CreatedAt: comment.CreatedAt,
		UpdatedAt: comment.UpdatedAt,
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Reply created successfully",
		"comment": commentResponse,
	})
}

// GetComments handles getting all comments for a post
func GetComments(c *gin.Context) {
	postID := c.Param("id")

	// Check if post exists
	var post models.Post
	if err := config.DB.First(&post, "id = ?", postID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	// Get comments with nested replies
	var comments []models.Comment
	if err := config.DB.Preload("Author").
		Preload("Replies", func(db *gorm.DB) *gorm.DB {
			return db.Preload("Author")
		}).
		Where("post_id = ? AND parent_comment_id IS NULL", postID).
		Order("created_at ASC").
		Find(&comments).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch comments"})
		return
	}

	// Convert to response format
	var commentsResponse []models.CommentResponse
	for _, comment := range comments {
		commentResponse := convertCommentToResponse(comment)
		commentsResponse = append(commentsResponse, commentResponse)
	}

	c.JSON(http.StatusOK, gin.H{
		"comments": commentsResponse,
	})
}

// UpdateComment handles updating a comment
func UpdateComment(c *gin.Context) {
	commentID := c.Param("id")

	// Get user from context
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	userModel := user.(models.User)

	// Find comment
	var comment models.Comment
	if err := config.DB.First(&comment, "id = ?", commentID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Comment not found"})
		return
	}

	// Check if user is the author
	if comment.AuthorID != userModel.ID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You can only update your own comments"})
		return
	}

	// Parse update request
	var req models.CommentUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update comment
	comment.Content = req.Content
	if err := config.DB.Save(&comment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update comment"})
		return
	}

	// Load author information
	config.DB.Preload("Author").First(&comment, comment.ID)

	commentResponse := models.CommentResponse{
		ID:              comment.ID,
		PostID:          comment.PostID,
		AuthorID:        comment.AuthorID,
		ParentCommentID: comment.ParentCommentID,
		Content:         comment.Content,
		Author: models.UserResponse{
			ID:        comment.Author.ID,
			Username:  comment.Author.Username,
			Email:     comment.Author.Email,
			CreatedAt: comment.Author.CreatedAt,
		},
		CreatedAt: comment.CreatedAt,
		UpdatedAt: comment.UpdatedAt,
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Comment updated successfully",
		"comment": commentResponse,
	})
}

// DeleteComment handles deleting a comment
func DeleteComment(c *gin.Context) {
	commentID := c.Param("id")

	// Get user from context
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	userModel := user.(models.User)

	// Find comment
	var comment models.Comment
	if err := config.DB.First(&comment, "id = ?", commentID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Comment not found"})
		return
	}

	// Check if user is the author
	if comment.AuthorID != userModel.ID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You can only delete your own comments"})
		return
	}

	// Delete comment (cascade will handle replies)
	if err := config.DB.Delete(&comment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete comment"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Comment deleted successfully",
	})
}
