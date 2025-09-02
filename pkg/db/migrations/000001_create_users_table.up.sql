CREATE TABLE IF NOT EXISTS users (
    user_id        UUID PRIMARY KEY UNIQUE ,
    first_name       VARCHAR(255) NOT NULL ,
    last_name       VARCHAR(255) NOT NULL
);