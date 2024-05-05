CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY,
    username VARCHAR(255) NOT NULL UNIQUE,
    email VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    profile_picture VARCHAR(255)
);

CREATE TABLE IF NOT EXISTS chats (
    id UUID PRIMARY KEY,
    message TEXT NOT NULL,
    room VARCHAR(255) DEFAULT 'general',
    from_user UUID NOT NULL,
    sent TIME DEFAULT CURRENT_TIME,
    FOREIGN KEY (from_user) REFERENCES users(id)
);

CREATE INDEX chatroom ON chats (room);
