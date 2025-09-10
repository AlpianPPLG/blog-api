package handlers

import (
	"net/http"
	"strconv"

	"blog-api/config"
	"blog-api/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// CreatePost handles creating a new post
func CreatePost(c *gin.Context) {
	var req models.PostCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user from context
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	userModel := user.(models.User)

	// Create post
	post := models.Post{
		ID:       uuid.New().String(),
		Title:    req.Title,
		Content:  req.Content,
		Tags:     req.Tags,
		AuthorID: userModel.ID,
	}

	if err := config.DB.Create(&post).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create post"})
		return
	}

	// Load author information
	config.DB.Preload("Author").First(&post, post.ID)

	// Convert to response format
	postResponse := models.PostResponse{
		ID:        post.ID,
		Title:     post.Title,
		Content:   post.Content,
		Tags:      post.Tags,
		AuthorID:  post.AuthorID,
		Author: models.UserResponse{
			ID:        post.Author.ID,
			Username:  post.Author.Username,
			Email:     post.Author.Email,
			CreatedAt: post.Author.CreatedAt,
		},
		LikesCount: 0,
		CreatedAt:  post.CreatedAt,
		UpdatedAt:  post.UpdatedAt,
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Post created successfully",
		"post":    postResponse,
	})
}

// GetPosts handles getting all posts with pagination
func GetPosts(c *gin.Context) {
	// Get pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset := (page - 1) * limit

	// Get posts with author and likes count
	var posts []models.Post
	var total int64

	// Count total posts
	config.DB.Model(&models.Post{}).Count(&total)

	// Get posts with pagination
	if err := config.DB.Preload("Author").
		Preload("Likes").
		Offset(offset).
		Limit(limit).
		Order("created_at DESC").
		Find(&posts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch posts"})
		return
	}

	// Convert to response format
	var postsResponse []models.PostResponse
	for _, post := range posts {
		postResponse := models.PostResponse{
			ID:        post.ID,
			Title:     post.Title,
			Content:   post.Content,
			Tags:      post.Tags,
			AuthorID:  post.AuthorID,
			Author: models.UserResponse{
				ID:        post.Author.ID,
				Username:  post.Author.Username,
				Email:     post.Author.Email,
				CreatedAt: post.Author.CreatedAt,
			},
			LikesCount: len(post.Likes),
			CreatedAt:  post.CreatedAt,
			UpdatedAt:  post.UpdatedAt,
		}
		postsResponse = append(postsResponse, postResponse)
	}

	c.JSON(http.StatusOK, gin.H{
		"posts": postsResponse,
		"pagination": gin.H{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	})
}

// GetPost handles getting a single post by ID
func GetPost(c *gin.Context) {
	postID := c.Param("id")

	var post models.Post
	if err := config.DB.Preload("Author").
		Preload("Comments", func(db *gorm.DB) *gorm.DB {
			return db.Preload("Author").Preload("Replies", func(db *gorm.DB) *gorm.DB {
				return db.Preload("Author")
			}).Where("parent_comment_id IS NULL")
		}).
		Preload("Likes").
		First(&post, "id = ?", postID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	// Convert comments to response format
	var commentsResponse []models.CommentResponse
	for _, comment := range post.Comments {
		commentResponse := convertCommentToResponse(comment)
		commentsResponse = append(commentsResponse, commentResponse)
	}

	// Convert likes to response format
	var likesResponse []models.LikeResponse
	for _, like := range post.Likes {
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

	postResponse := models.PostResponse{
		ID:        post.ID,
		Title:     post.Title,
		Content:   post.Content,
		Tags:      post.Tags,
		AuthorID:  post.AuthorID,
		Author: models.UserResponse{
			ID:        post.Author.ID,
			Username:  post.Author.Username,
			Email:     post.Author.Email,
			CreatedAt: post.Author.CreatedAt,
		},
		Comments:   commentsResponse,
		Likes:      likesResponse,
		LikesCount: len(post.Likes),
		CreatedAt:  post.CreatedAt,
		UpdatedAt:  post.UpdatedAt,
	}

	c.JSON(http.StatusOK, gin.H{
		"post": postResponse,
	})
}

// UpdatePost handles updating a post
func UpdatePost(c *gin.Context) {
	postID := c.Param("id")

	// Get user from context
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	userModel := user.(models.User)

	// Find post
	var post models.Post
	if err := config.DB.First(&post, "id = ?", postID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	// Check if user is the author
	if post.AuthorID != userModel.ID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You can only update your own posts"})
		return
	}

	// Parse update request
	var req models.PostUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update post
	updates := make(map[string]interface{})
	if req.Title != "" {
		updates["title"] = req.Title
	}
	if req.Content != "" {
		updates["content"] = req.Content
	}
	if req.Tags != nil {
		updates["tags"] = req.Tags
	}

	if err := config.DB.Model(&post).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update post"})
		return
	}

	// Reload post with author
	config.DB.Preload("Author").First(&post, post.ID)

	postResponse := models.PostResponse{
		ID:        post.ID,
		Title:     post.Title,
		Content:   post.Content,
		Tags:      post.Tags,
		AuthorID:  post.AuthorID,
		Author: models.UserResponse{
			ID:        post.Author.ID,
			Username:  post.Author.Username,
			Email:     post.Author.Email,
			CreatedAt: post.Author.CreatedAt,
		},
		LikesCount: 0,
		CreatedAt:  post.CreatedAt,
		UpdatedAt:  post.UpdatedAt,
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Post updated successfully",
		"post":    postResponse,
	})
}

// DeletePost handles deleting a post
func DeletePost(c *gin.Context) {
	postID := c.Param("id")

	// Get user from context
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	userModel := user.(models.User)

	// Find post
	var post models.Post
	if err := config.DB.First(&post, "id = ?", postID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	// Check if user is the author
	if post.AuthorID != userModel.ID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You can only delete your own posts"})
		return
	}

	// Delete post (cascade will handle comments and likes)
	if err := config.DB.Delete(&post).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete post"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Post deleted successfully",
	})
}

// convertCommentToResponse converts a comment to response format
func convertCommentToResponse(comment models.Comment) models.CommentResponse {
	// Convert replies
	var repliesResponse []models.CommentResponse
	for _, reply := range comment.Replies {
		replyResponse := models.CommentResponse{
			ID:              reply.ID,
			PostID:          reply.PostID,
			AuthorID:        reply.AuthorID,
			ParentCommentID: reply.ParentCommentID,
			Content:         reply.Content,
			Author: models.UserResponse{
				ID:        reply.Author.ID,
				Username:  reply.Author.Username,
				Email:     reply.Author.Email,
				CreatedAt: reply.Author.CreatedAt,
			},
			CreatedAt: reply.CreatedAt,
			UpdatedAt: reply.UpdatedAt,
		}
		repliesResponse = append(repliesResponse, replyResponse)
	}

	return models.CommentResponse{
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
		Replies:   repliesResponse,
		CreatedAt: comment.CreatedAt,
		UpdatedAt: comment.UpdatedAt,
	}
}
