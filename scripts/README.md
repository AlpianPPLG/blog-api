# Database Scripts for Blog API

This folder contains SQL scripts to set up the MySQL database for the Blog API.

## Scripts Overview

1. **01_create_database.sql** - Creates the `blog_api` database
2. **02_create_tables.sql** - Creates all required tables (users, posts, comments, likes)
3. **03_create_indexes.sql** - Creates indexes for better performance
4. **04_sample_data.sql** - Inserts sample data for testing

## How to Run

Execute the scripts in order using MySQL command line or any MySQL client:

```bash
mysql -u your_username -p < scripts/01_create_database.sql
mysql -u your_username -p < scripts/02_create_tables.sql
mysql -u your_username -p < scripts/03_create_indexes.sql
mysql -u your_username -p < scripts/04_sample_data.sql
```

Or run all at once:

```bash
mysql -u your_username -p < scripts/01_create_database.sql && \
mysql -u your_username -p < scripts/02_create_tables.sql && \
mysql -u your_username -p < scripts/03_create_indexes.sql && \
mysql -u your_username -p < scripts/04_sample_data.sql
```

## Database Schema

- **users**: Stores user information for authentication
- **posts**: Stores blog posts with tags (JSON format)
- **comments**: Stores comments with support for nested replies
- **likes**: Stores user likes for posts (many-to-many relationship)

## Sample Data

The sample data includes:
- 3 test users (password for all: "password")
- 3 sample blog posts
- Sample comments with nested replies
- Sample likes

## Environment Variables

Make sure to set these environment variables in your `.env` file:

```
DB_HOST=localhost
DB_PORT=3306
DB_USER=your_username
DB_PASSWORD=your_password
DB_NAME=blog_api
JWT_SECRET=your_jwt_secret_key
```
