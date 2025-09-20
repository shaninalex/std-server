CREATE TABLE users
(
    id            TEXT PRIMARY KEY NOT NULL,
    name          TEXT             NOT NULL,
    email         TEXT             NOT NULL UNIQUE,
    password_hash TEXT             NOT NULL,
    active        BOOLEAN          NOT NULL DEFAULT 0,
    created_at    DATETIME         NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at    DATETIME         NOT NULL DEFAULT CURRENT_TIMESTAMP
);


CREATE TABLE user_sessions
(
    id         TEXT PRIMARY KEY NOT NULL,
    user_id    TEXT             NOT NULL,
    data       TEXT             NOT NULL,
    expires_at DATETIME         NOT NULL,
    created_at DATETIME         NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);