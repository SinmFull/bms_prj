CREATE TABLE sensor_devices (
    id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
    building_id BIGINT UNSIGNED NOT NULL,
    sensor_type_id BIGINT UNSIGNED NOT NULL,
    name VARCHAR(255) NOT NULL COMMENT 'Name of sensor ,such as temperature#1',
    location TEXT COMMENT 'location of the sensor',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (building_id) REFERENCES buildings(id) ON DELETE CASCADE,
    FOREIGN KEY (sensor_type_id) REFERENCES sensor_types(id) ON DELETE CASCADE
);