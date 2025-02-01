CREATE TABLE users (
                       id SERIAL PRIMARY KEY,
                       userid VARCHAR(255) NOT NULL UNIQUE, -- Add UNIQUE constraint on userid
                       username VARCHAR(255) NOT NULL,
                       registered BOOLEAN DEFAULT FALSE,
                       language VARCHAR(3) DEFAULT 'ru',
                       timer TIME NOT NULL DEFAULT '22:00',
                       state VARCHAR NULL,
                       created_at TIMESTAMP DEFAULT (NOW() AT TIME ZONE 'Asia/Almaty')
);


CREATE TABLE reading_logs (
                              id SERIAL PRIMARY KEY,
                              userid VARCHAR(255) NOT NULL REFERENCES users(userid) ON DELETE CASCADE, -- Match the type with users table
                              date DATE NOT NULL DEFAULT CURRENT_DATE,
                              minutes_read INT NOT NULL CHECK (minutes_read >= 0),
                              created_at TIMESTAMP DEFAULT (NOW() AT TIME ZONE 'Asia/Almaty'),
                              UNIQUE (userid, date) -- Ensure unique constraint on the combination of userid and date
);

