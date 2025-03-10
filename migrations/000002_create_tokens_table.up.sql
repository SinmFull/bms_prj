CREATE TABLE IF NOT EXISTS tokens (
    hash BLOB,
    user_id BIGINT UNSIGNED NOT NULL,
    expiry DATETIME NOT NULL,
    scope TEXT NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
