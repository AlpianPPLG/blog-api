# Blog API with Golang

A RESTful API for a blog platform built with Go, featuring user authentication, post management, nested comments, and like functionality.

## Features

- **User Authentication**: JWT-based authentication with registration and login
- **Posts Management**: Full CRUD operations for blog posts
- **Nested Comments**: Support for comments and replies with hierarchical structure
- **Like System**: Users can like/unlike posts
- **Rate Limiting**: Protection against spam and abuse
- **Logging**: Comprehensive request/response logging
- **Database**: MySQL with GORM ORM

## Tech Stack

- **Language**: Go 1.21
- **Framework**: Gin
- **Database**: MySQL
- **ORM**: GORM
- **Authentication**: JWT
- **Logging**: Zap
- **Rate Limiting**: golang.org/x/time/rate

## Project Structure

```
blog-api/
├── config/
│   └── database.go          # Database configuration
├── handlers/
│   ├── auth.go              # Authentication handlers
│   ├── posts.go             # Post CRUD handlers
│   ├── comments.go          # Comment CRUD handlers
│   └── likes.go             # Like/unlike handlers
├── middleware/
│   ├── auth.go              # JWT authentication middleware
│   ├── rate_limit.go        # Rate limiting middleware
│   └── logging.go           # Logging middleware
├── models/
│   ├── user.go              # User model
│   ├── post.go              # Post model
│   ├── comment.go           # Comment model
│   └── like.go              # Like model
├── routes/
│   └── routes.go            # Route configuration
├── scripts/
│   ├── 01_create_database.sql
│   ├── 02_create_tables.sql
│   ├── 03_create_indexes.sql
│   ├── 04_sample_data.sql
│   └── README.md
├── go.mod
├── go.sum
├── main.go
├── env.example
└── README.md
```

## Setup Instructions

### 1. Prerequisites

- Go 1.21 or higher
- MySQL 5.7 or higher
- Git

### 2. Database Setup

1. Create a MySQL database and user
2. Run the SQL scripts in order:

```bash
mysql -u your_username -p < scripts/01_create_database.sql
mysql -u your_username -p < scripts/02_create_tables.sql
mysql -u your_username -p < scripts/03_create_indexes.sql
mysql -u your_username -p < scripts/04_sample_data.sql
```

### 3. Environment Configuration

1. Copy `env.example` to `.env`
2. Update the environment variables:

```env
DB_HOST=localhost
DB_PORT=3306
DB_USER=your_username
DB_PASSWORD=your_password
DB_NAME=blog_api
JWT_SECRET=your_jwt_secret_key_here
PORT=8080
```

### 4. Install Dependencies

```bash
go mod tidy
```

### 5. Run the Application

```bash
go run main.go
```

The API will be available at `http://localhost:8080`

## API Endpoints

### Authentication

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| POST | `/api/v1/auth/register` | Register new user | No |
| POST | `/api/v1/auth/login` | Login user | No |

### Posts

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| GET | `/api/v1/posts` | Get all posts (paginated) | No |
| GET | `/api/v1/posts/{id}` | Get post by ID | No |
| POST | `/api/v1/posts` | Create new post | Yes |
| PUT | `/api/v1/posts/{id}` | Update post | Yes (author only) |
| DELETE | `/api/v1/posts/{id}` | Delete post | Yes (author only) |

### Comments

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| GET | `/api/v1/posts/{id}/comments` | Get comments for post | No |
| POST | `/api/v1/posts/{id}/comments` | Create comment on post | Yes |
| POST | `/api/v1/comments/{id}/reply` | Reply to comment | Yes |
| PUT | `/api/v1/comments/{id}` | Update comment | Yes (author only) |
| DELETE | `/api/v1/comments/{id}` | Delete comment | Yes (author only) |

### Likes

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| POST | `/api/v1/posts/{id}/like` | Like a post | Yes |
| POST | `/api/v1/posts/{id}/unlike` | Unlike a post | Yes |
| GET | `/api/v1/posts/{id}/likes` | Get all likes for post | No |
| GET | `/api/v1/posts/{id}/like-status` | Check if user liked post | Yes |

### User Profile

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| GET | `/api/v1/profile` | Get current user profile | Yes |

## Request/Response Examples

### Register User

```bash
POST /api/v1/auth/register
Content-Type: application/json

{
  "username": "john_doe",
  "email": "john@example.com",
  "password": "password123"
}
```

### Create Post

```bash
POST /api/v1/posts
Authorization: Bearer <jwt_token>
Content-Type: application/json

{
  "title": "My First Post",
  "content": "This is the content of my first post.",
  "tags": ["golang", "api", "tutorial"]
}
```

### Create Comment

```bash
POST /api/v1/posts/{post_id}/comments
Authorization: Bearer <jwt_token>
Content-Type: application/json

{
  "content": "Great post! Thanks for sharing."
}
```

### Reply to Comment

```bash
POST /api/v1/comments/{comment_id}/reply
Authorization: Bearer <jwt_token>
Content-Type: application/json

{
  "content": "I agree with your point."
}
```

## Rate Limiting

- **Global Rate Limit**: 10 requests per second per IP
- **User Rate Limit**: 5 requests per second per authenticated user
- **Rate Limit Headers**: Included in responses when limit is exceeded

## Security Features

- JWT-based authentication
- Password hashing with bcrypt
- Rate limiting to prevent abuse
- Input validation
- SQL injection protection via GORM
- CORS support

## Testing

Use the provided Postman collection or test with curl:

```bash
# Health check
curl http://localhost:8080/health

# Register user
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","email":"test@example.com","password":"password123"}'

# Login
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123"}'
```

## Sample Data

The database includes sample data:
- 3 test users (password: "password")
- 3 sample blog posts
- Sample comments with nested replies
- Sample likes

## Error Handling

The API returns appropriate HTTP status codes and error messages:

- `200 OK`: Success
- `201 Created`: Resource created successfully
- `400 Bad Request`: Invalid request data
- `401 Unauthorized`: Authentication required
- `403 Forbidden`: Access denied
- `404 Not Found`: Resource not found
- `409 Conflict`: Resource already exists
- `429 Too Many Requests`: Rate limit exceeded
- `500 Internal Server Error`: Server error

## Development

### Running in Development Mode

```bash
# Set environment
export GIN_MODE=debug

# Run with hot reload (requires air)
air

# Or run directly
go run main.go
```

### Database Migrations

The application uses GORM's AutoMigrate feature. Models are automatically migrated on startup.

### Logging

Logs are structured using Zap logger with different levels:
- Info: General application flow
- Error: Error conditions
- Fatal: Critical errors that cause shutdown

## License

This project is open source and available under the MIT License.
