package routes

import (
	"blog-api/handlers"
	"blog-api/middleware"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// SetupRoutes configures all routes for the API
func SetupRoutes(logger interface{}) *gin.Engine {
	// Set Gin mode
	gin.SetMode(gin.ReleaseMode)

	// Create Gin router
	r := gin.New()

	// Add middleware
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(middleware.RateLimitMiddleware(rate.Every(1), 10)) // 10 requests per second

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "Blog API is running",
		})
	})

	// API v1 routes
	v1 := r.Group("/api/v1")
	{
		// Auth routes (no authentication required)
		auth := v1.Group("/auth")
		{
			auth.POST("/register", handlers.Register)
			auth.POST("/login", handlers.Login)
		}

		// Public routes (no authentication required)
		public := v1.Group("/")
		{
			// Posts (public read access)
			public.GET("/posts", handlers.GetPosts)
			public.GET("/posts/:id", handlers.GetPost)
			public.GET("/posts/:id/comments", handlers.GetComments)
			public.GET("/posts/:id/likes", handlers.GetPostLikes)
		}

		// Protected routes (authentication required)
		protected := v1.Group("/")
		protected.Use(middleware.AuthMiddleware())
		protected.Use(middleware.UserRateLimitMiddleware(rate.Every(1), 5)) // 5 requests per second per user
		{
			// User profile
			protected.GET("/profile", handlers.GetProfile)

			// Posts (authenticated)
			protected.POST("/posts", handlers.CreatePost)
			protected.PUT("/posts/:id", handlers.UpdatePost)
			protected.DELETE("/posts/:id", handlers.DeletePost)

			// Comments (authenticated)
			protected.POST("/posts/:id/comments", handlers.CreateComment)
			protected.PUT("/comments/:id", handlers.UpdateComment)
			protected.DELETE("/comments/:id", handlers.DeleteComment)

			// Replies (authenticated)
			protected.POST("/comments/:id/reply", handlers.ReplyToComment)

			// Likes (authenticated)
			protected.POST("/posts/:id/like", handlers.LikePost)
			protected.POST("/posts/:id/unlike", handlers.UnlikePost)
			protected.GET("/posts/:id/like-status", handlers.CheckUserLike)
		}
	}

	return r
}
