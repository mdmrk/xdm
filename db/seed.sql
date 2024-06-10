-- Insert fake data into the 'users' table
DO $$
BEGIN
    FOR i IN 1..100 LOOP
        INSERT INTO users (alias, username, password, salt, token)
        VALUES (
            'User' || i,
            'user' || lpad(i::text, 3, '0'),
            crypt('password' || i, gen_salt('bf')),
            gen_random_bytes(16),
            gen_random_bytes(16)
        )
        ON CONFLICT DO NOTHING;
    END LOOP;
END $$;

-- Insert fake data into the 'posts' table
DO $$
DECLARE
    user_ids UUID[];
BEGIN
    SELECT array_agg(id) INTO user_ids FROM users;

    FOR i IN 1..500 LOOP
        INSERT INTO posts (user_id, body)
        VALUES (
            user_ids[floor(random() * array_length(user_ids, 1) + 1)::int],
            'This is post number ' || i
        )
        ON CONFLICT DO NOTHING;
    END LOOP;
END $$;

-- Insert fake data into the 'follows' table
DO $$
DECLARE
    user_ids UUID[];
BEGIN
    SELECT array_agg(id) INTO user_ids FROM users;

    FOR i IN 1..2000 LOOP
        INSERT INTO follows (follower_id, followed_id)
        VALUES (
            user_ids[floor(random() * array_length(user_ids, 1) + 1)::int],
            user_ids[floor(random() * array_length(user_ids, 1) + 1)::int]
        )
        ON CONFLICT DO NOTHING;
    END LOOP;
END $$;

-- Insert fake data into the 'likes' table
DO $$
DECLARE
    user_ids UUID[];
    post_ids UUID[];
BEGIN
    SELECT array_agg(id) INTO user_ids FROM users;
    SELECT array_agg(id) INTO post_ids FROM posts;

    FOR i IN 1..20000 LOOP
        INSERT INTO likes (post_id, user_id)
        VALUES (
            post_ids[floor(random() * array_length(post_ids, 1) + 1)::int],
            user_ids[floor(random() * array_length(user_ids, 1) + 1)::int]
        )
        ON CONFLICT DO NOTHING;
    END LOOP;
END $$;

