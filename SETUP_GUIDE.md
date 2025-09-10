# Blog API Setup Guide

## Quick Start

### 1. Prerequisites
- Go 1.21 or higher
- MySQL 5.7 or higher
- Git

### 2. Database Setup

1. **Create MySQL database:**
```sql
CREATE DATABASE blog_api;
```

2. **Run the SQL scripts in order:**
```bash
# On Windows (Command Prompt)
mysql -u root -p < scripts\01_create_database.sql
mysql -u root -p < scripts\02_create_tables.sql
mysql -u root -p < scripts\03_create_indexes.sql
mysql -u root -p < scripts\04_sample_data.sql

# On Linux/Mac
mysql -u root -p < scripts/01_create_database.sql
mysql -u root -p < scripts/02_create_tables.sql
mysql -u root -p < scripts/03_create_indexes.sql
mysql -u root -p < scripts/04_sample_data.sql
```

### 3. Environment Configuration

1. **Copy environment file:**
```bash
copy env.example .env
```

2. **Edit .env file with your database credentials:**
```env
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
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

## Testing the API

### Using Postman
1. Import the `postman_collection.json` file
2. Set the `base_url` variable to `http://localhost:8080`
3. Start with the "Register User" request
4. Copy the JWT token from the response
5. Set the `jwt_token` variable with the token
6. Test other endpoints

### Using curl (Linux/Mac)
```bash
chmod +x test_api.sh
./test_api.sh
```

### Using curl (Windows)
```cmd
test_api.bat
```

## API Endpoints Summary

### Authentication
- `POST /api/v1/auth/register` - Register new user
- `POST /api/v1/auth/login` - Login user
- `GET /api/v1/profile` - Get user profile (requires auth)

### Posts
- `GET /api/v1/posts` - Get all posts (public)
- `GET /api/v1/posts/{id}` - Get post by ID (public)
- `POST /api/v1/posts` - Create post (requires auth)
- `PUT /api/v1/posts/{id}` - Update post (requires auth, author only)
- `DELETE /api/v1/posts/{id}` - Delete post (requires auth, author only)

### Comments
- `GET /api/v1/posts/{id}/comments` - Get comments for post (public)
- `POST /api/v1/posts/{id}/comments` - Create comment (requires auth)
- `POST /api/v1/comments/{id}/reply` - Reply to comment (requires auth)
- `PUT /api/v1/comments/{id}` - Update comment (requires auth, author only)
- `DELETE /api/v1/comments/{id}` - Delete comment (requires auth, author only)

### Likes
- `POST /api/v1/posts/{id}/like` - Like post (requires auth)
- `POST /api/v1/posts/{id}/unlike` - Unlike post (requires auth)
- `GET /api/v1/posts/{id}/likes` - Get all likes for post (public)
- `GET /api/v1/posts/{id}/like-status` - Check if user liked post (requires auth)

## Sample Data

The database includes sample data:
- **Users**: john_doe, jane_smith, bob_wilson (password: "password")
- **Posts**: 3 sample blog posts with different topics
- **Comments**: Sample comments with nested replies
- **Likes**: Sample likes on posts

## Troubleshooting

### Common Issues

1. **Database Connection Error**
   - Check MySQL is running
   - Verify database credentials in .env file
   - Ensure database exists

2. **Port Already in Use**
   - Change PORT in .env file
   - Kill process using port 8080

3. **JWT Token Issues**
   - Ensure JWT_SECRET is set in .env
   - Check token format in Authorization header

4. **Rate Limiting**
   - Global: 10 requests/second per IP
   - User: 5 requests/second per authenticated user

### Logs
The application logs all requests and errors. Check the console output for detailed information.

## Project Structure

```
blog-api/
├── config/          # Database configuration
├── handlers/        # HTTP handlers
├── middleware/      # Middleware functions
├── models/          # Data models
├── routes/          # Route definitions
├── scripts/         # SQL scripts
├── main.go          # Application entry point
└── README.md        # Documentation
```

## Features Implemented

✅ User authentication with JWT
✅ Posts CRUD operations
✅ Nested comments system
✅ Like/unlike functionality
✅ Rate limiting
✅ Request logging
✅ Input validation
✅ Error handling
✅ Database relationships
✅ Pagination support
✅ Security features

## Next Steps

1. Test all endpoints using Postman
2. Verify database relationships
3. Test rate limiting
4. Check error handling
5. Monitor logs for any issues

The API is now ready for use!
