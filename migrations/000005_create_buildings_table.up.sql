CREATE TABLE buildings (
    id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
    group_id BIGINT UNSIGNED NOT NULL,
    name VARCHAR(255) NOT NULL,
    location TEXT COMMENT 'location of the building',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (group_id) REFERENCES user_groups(id) ON DELETE CASCADE
);