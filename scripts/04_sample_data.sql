-- Insert sample data for testing
USE blog_api;

-- Sample users
INSERT INTO users (id, username, email, password_hash) VALUES
('550e8400-e29b-41d4-a716-446655440001', 'john_doe', 'john@example.com', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi'),
('550e8400-e29b-41d4-a716-446655440002', 'jane_smith', 'jane@example.com', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi'),
('550e8400-e29b-41d4-a716-446655440003', 'bob_wilson', 'bob@example.com', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi');

-- Sample posts
INSERT INTO posts (id, title, content, tags, author_id) VALUES
('660e8400-e29b-41d4-a716-446655440001', 'Getting Started with Go', 'Go is a programming language developed by Google...', '["golang", "programming", "tutorial"]', '550e8400-e29b-41d4-a716-446655440001'),
('660e8400-e29b-41d4-a716-446655440002', 'Building REST APIs', 'REST APIs are a way to provide web services...', '["api", "rest", "web"]', '550e8400-e29b-41d4-a716-446655440002'),
('660e8400-e29b-41d4-a716-446655440003', 'Database Design Best Practices', 'When designing databases, consider normalization...', '["database", "design", "sql"]', '550e8400-e29b-41d4-a716-446655440001');

-- Sample comments
INSERT INTO comments (id, post_id, author_id, parent_comment_id, content) VALUES
('770e8400-e29b-41d4-a716-446655440001', '660e8400-e29b-41d4-a716-446655440001', '550e8400-e29b-41d4-a716-446655440002', NULL, 'Great tutorial! Very helpful for beginners.'),
('770e8400-e29b-41d4-a716-446655440002', '660e8400-e29b-41d4-a716-446655440001', '550e8400-e29b-41d4-a716-446655440003', '770e8400-e29b-41d4-a716-446655440001', 'I agree! This helped me understand Go better.'),
('770e8400-e29b-41d4-a716-446655440003', '660e8400-e29b-41d4-a716-446655440002', '550e8400-e29b-41d4-a716-446655440001', NULL, 'Nice explanation of REST principles.');

-- Sample likes
INSERT INTO likes (id, post_id, user_id) VALUES
('880e8400-e29b-41d4-a716-446655440001', '660e8400-e29b-41d4-a716-446655440001', '550e8400-e29b-41d4-a716-446655440002'),
('880e8400-e29b-41d4-a716-446655440002', '660e8400-e29b-41d4-a716-446655440001', '550e8400-e29b-41d4-a716-446655440003'),
('880e8400-e29b-41d4-a716-446655440003', '660e8400-e29b-41d4-a716-446655440002', '550e8400-e29b-41d4-a716-446655440001');
