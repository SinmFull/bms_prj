CREATE TABLE  sensor_value_day(
    id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
    sensor_device_id VARCHAR(255) NOT NULL,
    value text NOT NULL COMMENT 'Sensor daily reading',
    timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);