CREATE TABLE  sensor_value(
    id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
    sensor_device_id BIGINT UNSIGNED NOT NULL,
    value text NOT NULL COMMENT 'Sensor reading',
    timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (sensor_device_id) REFERENCES sensor_devices(id) ON DELETE CASCADE
);