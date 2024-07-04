SET session_replication_role = 'replica';

TRUNCATE TABLE messages CASCADE;
TRUNCATE TABLE likes CASCADE;
TRUNCATE TABLE follows CASCADE;
TRUNCATE TABLE posts CASCADE;
TRUNCATE TABLE users CASCADE;

SET session_replication_role = 'origin';
