CREATE TABLE sensor_types (
    id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(100) UNIQUE NOT NULL COMMENT 'Sensor type, such as temperature, humidity, etc.',
    unit VARCHAR(50) NOT NULL COMMENT 'Unit, sucha as °C, %RH, ppm etc.'
);