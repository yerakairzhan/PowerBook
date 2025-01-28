CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    userid VARCHAR(255) NOT NULL,
    username VARCHAR(255) NOT NULL,
    registered BOOLEAN DEFAULT FALSE,
    language varchar(3) DEFAULT 'ru',
    created_at TIMESTAMP DEFAULT (NOW() AT TIME ZONE 'Asia/Almaty')
);
