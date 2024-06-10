CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    alias VARCHAR(32) NOT NULL,
    username VARCHAR(16) NOT NULL UNIQUE,
    password VARCHAR(64) NOT NULL,
    salt BYTEA NOT NULL,
    token BYTEA,
    seen TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    public_key TEXT,
    private_key TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE IF NOT EXISTS posts (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  user_id UUID NOT NULL,
  likes int DEFAULT 0,
  body VARCHAR(1024) NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  CONSTRAINT fk_user
        FOREIGN KEY (user_id)
        REFERENCES users (id)
        ON DELETE CASCADE
);
CREATE TABLE IF NOT EXISTS follows (
    follower_id UUID NOT NULL,
    followed_id UUID NOT NULL,
    follow_date TIMESTAMP WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (follower_id, followed_id),
    CONSTRAINT fk_follower
        FOREIGN KEY (follower_id)
        REFERENCES users (id)
        ON DELETE CASCADE,
    CONSTRAINT fk_followed
        FOREIGN KEY (followed_id)
        REFERENCES users (id)
        ON DELETE CASCADE
);
CREATE TABLE IF NOT EXISTS likes (
    post_id UUID NOT NULL,
    user_id UUID NOT NULL,
    liked_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (post_id, user_id),
    CONSTRAINT fk_post
        FOREIGN KEY (post_id)
        REFERENCES posts (id)
        ON DELETE CASCADE,
    CONSTRAINT fk_liker
        FOREIGN KEY (user_id)
        REFERENCES users (id)
        ON DELETE CASCADE
);
CREATE TABLE IF NOT EXISTS messages (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    sender_id UUID NOT NULL,
    recipient_id UUID NOT NULL,
    content BYTEA NOT NULL,
    timestamp TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT sender_id
        FOREIGN KEY (sender_id)
        REFERENCES users (id)
        ON DELETE CASCADE,
    CONSTRAINT recipient_id
        FOREIGN KEY (recipient_id)
        REFERENCES users (id)
        ON DELETE CASCADE
);
CREATE OR REPLACE FUNCTION increment_like_count()
RETURNS TRIGGER AS $$
BEGIN
    UPDATE posts
    SET likes = likes + 1
    WHERE id = NEW.post_id;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION decrement_like_count()
RETURNS TRIGGER AS $$
BEGIN
    UPDATE posts
    SET likes = likes - 1
    WHERE id = OLD.post_id;
    RETURN OLD;
END;
$$ LANGUAGE plpgsql;
CREATE TRIGGER increment_likes_after_insert
AFTER INSERT ON likes
FOR EACH ROW
EXECUTE FUNCTION increment_like_count();

CREATE TRIGGER decrement_likes_after_delete
AFTER DELETE ON likes
FOR EACH ROW
EXECUTE FUNCTION decrement_like_count();

