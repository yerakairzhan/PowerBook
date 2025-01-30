DROP TABLE IF EXISTS users;
ALTER TABLE reading_logs DROP CONSTRAINT reading_logs_userid_fkey;
drop table if exists reading_logs;