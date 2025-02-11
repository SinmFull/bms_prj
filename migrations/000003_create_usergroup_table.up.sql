CREATE TABLE user_groups (
    id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(255) NOT NULL UNIQUE COMMENT 'name of usergroup, default value is email of admin',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);